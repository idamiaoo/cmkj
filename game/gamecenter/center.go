package main

import (
	"time"

	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"
)

const (
	onlinePeriod = 1 * time.Second
)

type Center struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	size       int
}

var DefaultCenter *Center

func NewCenter() *Center {
	return &Center{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		size:       0,
	}
}

func (center *Center) AddClient(cli *Client) {
	util.Log.Debug("login game")
	center.clients[cli.player.UserName] = cli
}

func (center *Center) RemoveClient(cli *Client) {
	util.Log.Debug("exit game")
	if cli.player == nil {
		return
	}

	_, ok := center.clients[cli.player.UserName]
	if false == ok {
		return
	}
	delete(center.clients, cli.player.UserName)
}

//更新在线人数
func (center *Center) sendOnlineNum() {
	if len(center.clients) == center.size {
		return
	}
	center.size = len(center.clients)
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_ONLINENUM)
	buf.WriteInt(int32(len(center.clients)))
	msg := buf.Bytes()
	center.broadcast <- msg
	util.Log.Debug("sned online num:", center.size)
}

func (center *Center) run() {
	go func(center *Center) {
		trick := time.NewTicker(onlinePeriod)
		defer trick.Stop()
		for {
			select {
			case <-trick.C:
				center.sendOnlineNum()
			}
		}
	}(center)

	go func() {
		trick := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-trick.C:
				umembers := models.ReadMemberUpload()
				for _, umember := range umembers {
					getOutPlayer(umember.UserName)
				}
			}
		}
	}()

	for {
		util.Log.Debug("center run")
		select {
		case cli := <-center.register:
			center.AddClient(cli)
			util.Log.Debug("register over")
		case cli := <-center.unregister:
			center.RemoveClient(cli)
		case msg := <-center.broadcast:
			util.Log.Debug("brodacast")
			for _, cli := range center.clients {
				select {
				case cli.send <- msg:
				default:
					delete(center.clients, cli.getName())
				}
			}

		}
	}
}

func (center *Center) Start() {
	go center.run()
}
