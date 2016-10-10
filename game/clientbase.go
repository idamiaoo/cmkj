package game

import (
	//. "go/cmkj_server_go/game"
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"

	"net/http"
	//"strconv"
	//"strings"
	"time"

	//"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	reg     = ";"
	layout  = "2006-01-02 15:04:05"
)

//Upgrader  服务端websocket默认配置
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Client 游戏客户端
type Client struct {
	Conn *websocket.Conn
	P    *Player
	Send chan []byte
	Game IGame
	Max  int
}

//GetGame 获取客户端绑定的游戏
func (c *Client) GetGame() IGame {
	return c.Game
}

//MessageHandler 消息处理
func (c *Client) MessageHandler(client IClient, msg []byte) {
	buf := util.NewByteBufferWith(msg)
	head := buf.ReadShort()
	msgid := buf.ReadInt()
	switch head {
	case util.M_BET:
		util.Log.Debug("M_BET")
		num := buf.ReadByte()
		name := buf.ReadUTF()
		t := buf.ReadShort()
		index := buf.ReadShort()
		re := c.DoBet(name, int(t), int(index), int(num), buf.Bytes())
		util.Log.Debug(re)
		c.sendRe(head, msgid, re)
	case util.M_MG:
		util.Log.Debug("M_MG")
		num := buf.ReadByte()
		cmd := buf.ReadInt()
		args := buf.ReadUTF()
		arg := buf.ReadUTF()
		re := c.DoManCmd(int(num), int(cmd), args, arg)
		c.sendRe(head, msgid, re)
	case util.M_LOGIN:
		util.Log.Debug("M_LOGIN")
		num := buf.ReadByte()
		index := buf.ReadByte()
		t := buf.ReadByte()
		name := buf.ReadUTF()
		pwd := buf.ReadUTF()
		xy := buf.ReadUTF()
		ip := buf.ReadUTF()
		re := c.DoLogin(client, int(num), int(t), int(index), name, pwd, xy, ip)
		c.sendRe(head, msgid, re)
	case util.M_INROOM:
		util.Log.Debug("M_INROOM")
		num := buf.ReadByte()
		t := buf.ReadByte()
		re := c.InTable(int(num), int(t))
		c.sendRe(head, msgid, re)
	case util.M_EXIT:
		util.Log.Debug("M_EXIT")
		num := buf.ReadByte()
		re := c.LeaveTable(int(num))
		c.sendRe(head, msgid, re)
	case util.M_TIPS:
		util.Log.Debug("M_TIPS")
		num := buf.ReadByte()
		name := buf.ReadUTF()
		t := buf.ReadShort()
		total := buf.ReadDouble()
		re := c.DoTips(int(num), int(t), name, total)
		c.sendRe(head, msgid, re)
	case util.M_MGLOGIN:
		util.Log.Debug("M_MGLOGIN")
		num := buf.ReadByte()
		cmd := buf.ReadInt()
		args := buf.ReadUTF()
		arg := buf.ReadUTF()
		re := c.DoManLogin(int(num), int(cmd), args, arg)
		c.sendRe(head, msgid, re)
	case util.M_HEAR, util.M_HEAR1:
		//util.Log.Debug("M_HEAR")
		c.Hear()
		c.sendRe(head, msgid, 0)
	default:
		util.Log.Error("undefine cmd")
	}
}

func (c *Client) disconnectHandler() {
	return
}

//DoLogin 客户端登录(玩家)
func (c *Client) DoLogin(client IClient, num, t, platform int, name, pwd, xy, ip string) byte {

	if num < 1 || num > c.Max {
		util.Log.Errorf("num=%d\n", int(num))
		return 1
	}
	var isFirst bool
	p := c.P
	p.Name = name
	p.IP = ip
	p.Platform = platform
	p.Type = t
	p.Game = 1
	p.Xy = xy
	p.Send = c.Send
	member, re := models.ReadMember(p.Name, pwd)
	if re != 0 {
		return re
	}
	p.ID = member.UserID
	p.Balance = member.CurMoney1
	p.ClassID = member.ClassID
	p.Maxprofit = member.Maxprofit
	p.IsLock = member.IsLock
	p.Moneysort = member.MoneySort
	p.LoginTime = string(time.Now().Unix())
	p.PreSequence = member.PreSequence
	var ids string
	if p.Game == util.Roulette {
		ids = member.RouletteBetLimitIDs
	} else {
		ids = member.VideoBetLimitIDs
	}
	re = p.GetLimit(ids)
	if re != 0 {
		c.P = nil
		return re
	}
	re = models.UpdateLoginInfo(ip, name, c.P.Game)
	if re != 0 {
		c.P = nil
		return re
	}
	client.InTable(num, t)
	return 0
}

//InTable 进入牌桌(玩家)
func (c *Client) InTable(num, t int) byte {
	util.Log.Debug(num, t)
	if c.P == nil || num < 1 || num > c.Max {
		return 1
	}
	if c.P.Bets == nil {
		c.P.Bets = make(map[int]*Bet)
	}
	if _, ok := c.P.Bets[num]; !ok {
		c.P.Bets[num] = NewBet()
	}
	c.P.Online = true
	c.P.Home = num
	//发送首次消息
	go func() {
		c.Game.EnterRoom(c.P, t, num)
	}()
	return 0
}

//LeaveTable 离开牌桌(玩家)
func (c *Client) LeaveTable(num int) byte {
	if c.P == nil || num < 0 || num > c.Max {
		return 1
	}
	offline := false
	if num == 0 {
		num = c.P.Home
		offline = true
		c.P.Online = false
	}
	c.Game.LeaveRoom(c.P, offline, num)
	return 0
}

//DoManLogin 管理员登录(台面)
func (c *Client) DoManLogin(num, cmd int, args, arg string) byte {
	if num < 1 || num > c.Max {
		return 1
	}

	//房间管理登录处理
	return DoManageLogin(c.Game, c.P, c.Send, cmd, num, args, arg)
}

//DoTips 牌桌小费(玩家)
func (c *Client) DoTips(num, t int, name string, total float64) byte {
	if num < 1 || num > c.Max {
		return 1
	}
	return 0
}

//DoBet 下注(玩家)
func (c *Client) DoBet(name string, t, index, num int, body []byte) byte {
	if num < 1 || num > c.Max {
		return 1
	}
	//房间投注处理
	return c.Game.DoBet(c.P, name, t, index, num, body)
}

//DoManCmd 管理员命令(台面)
func (c *Client) DoManCmd(num, cmd int, args, arg string) byte {
	if num < 1 || num > c.Max {
		return 1
	}
	//房间管理命令处理
	return DoManageCmd(c.Game, c.P, c.Send, cmd, num, args, arg)
}

//sendRe 客户端消息响应
func (c *Client) sendRe(head int16, id int32, re byte) {
	buf := util.NewByteBuffer()
	buf.WriteShort(head)
	buf.WriteInt(id)
	buf.WriteByte(byte(re))
	msg := buf.Bytes()
	c.Send <- msg
	return
}

//ReadPump 客户端读go程
func (c *Client) ReadPump(client IClient) {
	defer c.disconnectHandler()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			util.Log.Warningf("%v", err)
			break
		}
		//默认是按二进制消息处理
		c.MessageHandler(client, message)
	}
}

func (c *Client) writeRaw(msg []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.Conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

//WritePump 客户端写go程
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message := <-c.Send
		if message == nil {
			// 客户端连接关闭
			return
		}
		if err := c.writeRaw(message); err != nil {
			util.Log.Errorf("%v", err)
			return
		}
		buf := util.NewByteBufferWith(message)
		cmd := int(buf.ReadShort())
		if cmd != 19999 && cmd != 29999 {
			util.Log.Debugf("[%s] send :%d", c.P.Name, cmd)
		}

	}
}

//Hear 心跳包处理
func (c *Client) Hear() {
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		util.Log.Error(err)
	}
}
