package network

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/util"

	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type HallClinet struct {
	Conn *websocket.Conn
	Send chan []byte
}

var HClient *HallClinet

func NewhallClient() *HallClinet {
	return &HallClinet{
		Send: make(chan []byte, 10),
	}
}

func (client *HallClinet) SendMsg(msg []byte) {
	client.Send <- msg
}

func (client *HallClinet) runOnce() {
	defer client.Conn.Close()
	done := make(chan []byte)
	go func() {
		defer close(done)
		for {
			_, message, err := client.Conn.ReadMessage()
			if err != nil {
				util.Log.Error(err)
				return
			}
			buf := util.NewByteBufferWith(message)
			cmd := buf.ReadShort()
			if cmd != 19999 {
				util.Log.Infof("recv: %d", cmd)
			}
		}
	}()
	tricker := time.NewTicker(time.Second * 10)
	defer tricker.Stop()
loop:
	for {
		select {
		case <-tricker.C: //心跳
			buf := util.NewByteBuffer()
			buf.WriteShort(util.M_HEAR)
			buf.WriteInt(1)
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
				util.Log.Error(err)
			}
		case msg := <-client.Send:
			if msg == nil {
				util.Log.Error("nil msg")
				break
			}
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				util.Log.Error(err)
			}
		case <-done:
			util.Log.Error("clinet close")
			break loop
		}
	}
	util.Log.Error("HallClient close")
}

func (client *HallClinet) run() {
	nsleep := 100
	u := url.URL{
		Scheme: "ws",
		Host:   conf.Conf.DefaultString("loginhost", "127.0.0.1:3000"),
		Path:   conf.Conf.DefaultString("loginpath", "/login"),
	}
	for {
		util.Log.Debugf("connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			util.Log.Errorf("connet login error: %s", err)
			nsleep *= 2
			if nsleep > 60*1000 {
				nsleep = 60 * 1000
			}
			util.Log.Infof("connect again after %d Millisecond ", nsleep)
			time.Sleep(time.Duration(nsleep) * time.Millisecond)
			continue
		}
		client.Conn = c
		util.Log.Info("connect login success")
		nsleep = 100
		client.runOnce()
	}
}

func (client *HallClinet) Start() {
	go client.run()
}
