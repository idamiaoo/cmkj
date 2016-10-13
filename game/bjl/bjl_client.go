package main

import (
	Game "go/cmkj_server_go/game"
	"go/cmkj_server_go/util"

	"github.com/gin-gonic/gin"
)

//IntSet ...
type IntSet map[int]struct{}

//NewIntSet ...
func NewIntSet() IntSet {
	return make(map[int]struct{})
}

//Add ...
func (set IntSet) Add(v int) {
	if _, ok := set[v]; ok {
		return
	}
	set[v] = struct{}{}
}

//BjlClient 百家乐游戏客户端
type BjlClient struct {
	Game.Client //匿名包含该game.Client
	rooms       IntSet
}

//LeaveTable 重写 LeaveTable函数
func (c *BjlClient) LeaveTable(num int) byte {
	if c.P == nil {
		return 1
	}
	closed := false
	if num <= 0 {
		if num == -9 {
			closed = false
		} else {
			closed = true
			c.P.Online = false
		}
		if c.P.Type < 0 {
			c.Game.LeaveRoom(c.P, closed, c.P.Home)
			return 1
		}
		if c.Game == nil {
			return 1
		}
		for tid := range c.rooms {
			c.Game.LeaveRoom(c.P, closed, tid)
		}
		c.rooms = NewIntSet()
	}
	if num < 1 || num > c.Max {
		return 1
	}
	c.Game.LeaveRoom(c.P, closed, num)
	delete(c.rooms, num)
	return 0
}

//InTable 重写InTable函数
func (c *BjlClient) InTable(num, t int) byte {
	util.Log.Debug(num, t)
	if c.P == nil || num < 1 || num > c.Max {
		return 1
	}
	if c.rooms == nil {
		c.rooms = NewIntSet()
	}
	if t != 2 && c.P.Home != num {
		c.LeaveTable(-9)
	}
	if c.P.Bets == nil {
		c.P.Bets = make(map[int]*Game.Bet)
	}
	if _, ok := c.P.Bets[num]; !ok {
		c.P.Bets[num] = Game.NewBet()
	}
	c.P.Online = true
	c.P.Home = num
	c.rooms.Add(num)
	go func() {
		c.Game.EnterRoom(c.P, t, num)
	}()
	util.Log.Debug("in table")
	return 0
}

func handleClient(c *gin.Context) {
	util.Log.Debug("aceept a cient !")
	conn, err := Game.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.Log.Error(err)
		return
	}
	client := &BjlClient{
		Client: Game.Client{
			Send: make(chan []byte, 256),
			Conn: conn,
			Game: Bjl,
			Max:  Bjl.Max,
		},
	}

	//client.LeaveTable()
	go client.WritePump()
	client.ReadPump(client)
}
