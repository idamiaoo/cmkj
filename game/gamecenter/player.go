package main

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/timer"
	"go/cmkj_server_go/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	//"bytes"
	//"fmt"
	"net/http"
	//"strconv"
	"strings"
	"time"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Player struct {
	UserName  string
	UserID    int64
	Blance    float64
	Type      byte
	Platform  byte
	Online    bool
	LimitMax  int
	LimitMin  int
	Chips     string
	Limits    string
	MoneySort string
	Nickname  string
	PreID     int64
	IsChat    byte
	IsTip     byte
}

func (p *Player) makeMsg() []byte {
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_PLAYER)
	buf.WriteUTF(p.UserName)
	buf.WriteDouble(p.Blance)
	buf.WriteUTF(p.Chips)
	buf.WriteUTF(p.MoneySort)
	buf.WriteByte(p.IsTip)
	buf.WriteByte(p.IsChat)
	buf.WriteUTF((p.Nickname))
	return buf.Bytes()
}

type Client struct {
	player *Player
	conn   *websocket.Conn
	send   chan []byte
	Games  map[string]struct{}
}

func (client *Client) getName() string {
	return client.player.UserName
}

func (client *Client) doLogin(num, t, index byte, name, pwd, way string) byte {
	if client.player == nil {
		client.player = &Player{}
	}
	p := client.player
	p.UserName = name
	p.Type = t
	p.Platform = index
	//验证登录
	if re := login(p, pwd); re != 0 {
		return re
	}
	//加载房间路由
	if strings.EqualFold("-1", way) {
		way = conf.Conf.String("games")
	}
	way = strings.TrimSuffix(way, ",")
	ways := strings.Split(way, ",")
	for _, v := range ways {
		c, _ := strconv.Atoi(v)
		if c > 0 {
			if _, ok := client.Games[v]; ok {
				continue
			}
			client.Games[v] = struct{}{}
		}
	}
	client.inRoom() //进入大厅
	go client.sendFirstMsg()
	return 0
}

func (client *Client) exitGame() byte {
	util.Log.Debug("exitGame")
	if client.player == nil {
		return 1
	}
	client.Games = make(map[string]struct{})
	c, _ := DefaultCenter.clients[client.player.UserName]
	if c != client {
		return 0
	}
	models.UpdateOlineOrExit(client.player.UserName, util.DO_EXIT)
	DefaultCenter.unregister <- client
	return 0
}

func (client *Client) inRoom() byte {
	if client.player == nil {
		return 1
	}
	c, ok := DefaultCenter.clients[client.player.UserName]
	if ok {
		if c != nil && c != client {
			c.getoutplayer(0)
		}
	}
	DefaultCenter.register <- client
	return 0
}

func (client *Client) changeNickName(name, pwd, nickname string) byte {
	util.Log.Debug("changeNickName")
	if client.player == nil {
		return 1
	}
	if !strings.EqualFold(name, client.player.UserName) {
		return 1
	}
	if !strings.EqualFold(nickname, "") {
		return 1
	}
	re := models.ChangeNickName(pwd, client.player.UserName, nickname)
	return re
}

func (client *Client) sendVideoURL() {
	p := client.player
	url := models.ReadVideoURL()
	util.Log.Debug("url:", url)
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_VIDEOURL)
	buf.WriteUTF(p.Limits)
	buf.WriteUTF(url)
	msg := buf.Bytes()
	client.send <- msg
}

func (client *Client) sendFirstMsg() {
	//发送玩家信息
	client.send <- client.player.makeMsg()
	//发送视频信息
	client.sendVideoURL()
	//发送游戏信息
	util.Log.Debug(client.Games)
	for way := range client.Games {
		gameID, tableID := getGame(way)
		status := &Status{
			GameType: gameID,
			TableID:  tableID,
		}
		if !getGamestatus(status) {
			continue
		}
		if -1 == status.Game {
			status = &Status{
				GameType: gameID,
				TableID:  tableID,
			}
		}
		//util.Log.Debug(status)
		util.Log.Debug(status.TableID)
		client.send <- status.makeMsg()
	}
}

func (client *Client) getoutplayer(now int32) {
	if client.player == nil {
		return
	}
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_GETOUT)
	buf.WriteInt(now)
	msg := buf.Bytes()
	if now == 0 {
		client.send <- msg
		client.send <- nil
		return
	}
	client.send <- msg
}

func (client *Client) messageHandler(msg []byte) {
	buf := util.NewByteBufferWith(msg)
	cmd := buf.ReadShort()
	switch cmd {
	case util.P_STATUS1:
		util.Log.Debug("status")
		//广播游戏状态
		g := buf.ReadByte()
		t := buf.ReadByte()
		game := int(g*10 + t)
		buf.Next(9)
		no := buf.ReadShort()
		util.Log.Debug(g, t, game, no)
		timer.GameTimer[game] = int(no)
		DefaultCenter.broadcast <- msg

	case util.M_LOGIN1, util.M_LOGIN2: //登录
		util.Log.Debug("login")
		id := buf.ReadInt()
		num := buf.ReadByte()
		index := buf.ReadByte()
		t := buf.ReadByte()
		name := strings.Replace(buf.ReadUTF(), reg, "", -1)
		pwd := buf.ReadUTF()
		way := buf.ReadUTF()
		util.Log.Debugf("pwd:[%s]\n", pwd)
		util.Log.Debugf("name:[%s]\n", name)
		util.Log.Debugf("way:[%s]\n", way)
		re := client.doLogin(num, t, index, name, pwd, way)
		client.sendReV1(cmd, id, re)
	case util.M_EXIT: //退出
		util.Log.Debug("exit!")
		id := buf.ReadInt()
		//num := buf.ReadByte()
		re := client.exitGame()
		client.sendRe(cmd, id, re)
	case util.M_HEAR, util.M_HEAR1: //心跳包
		//util.Log.Debug("heartbeat")
		client.hear()
		id := buf.ReadInt()
		client.sendRe(cmd, id, 0)
	case util.M_CHGNICK: //修改昵称
		name := buf.ReadUTF()
		pwd := buf.ReadUTF()
		nickname := buf.ReadUTF()
		re := client.changeNickName(name, pwd, nickname)
		client.sendRe(cmd, 0, re)
	default:
		util.Log.Debug("udefine data")
		//client.Close()
	}

}

func (client *Client) disconnectHandler() {
	DefaultCenter.unregister <- client
	client.send <- nil
}

func (client *Client) sendRe(head int16, id int32, re byte) {
	buf := util.NewByteBuffer()
	buf.WriteShort(head)
	buf.WriteInt(id)
	buf.WriteByte(re)
	msg := buf.Bytes()
	client.send <- msg
	return
}

func (client *Client) sendReV1(head int16, id int32, re byte) {
	buf := util.NewByteBuffer()
	buf.WriteShort(head)
	buf.WriteInt(id)
	buf.WriteByte(re)
	nowtime := time.Now().Format(layout)
	buf.WriteUTF(nowtime)
	msg := buf.Bytes()
	client.send <- msg
	return
}

func (client *Client) readPump() {
	defer client.disconnectHandler()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			util.Log.Warningf("%v", err)
			break
		}
		//默认是按二进制消息处理
		client.messageHandler(message)
	}
}

func (client *Client) writeRaw(msg []byte) error {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := client.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()
	for {
		message := <-client.send
		if message == nil {
			// 客户端连接关闭
			return
		}
		if err := client.writeRaw(message); err != nil {
			util.Log.Errorf("%v", err)
			return
		}
		buf := util.NewByteBufferWith(message)
		cmd := buf.ReadShort()
		if cmd != 29999 && cmd != 19999 {
			util.Log.Debug(cmd)
		}

	}
}

func (client *Client) hear() {
	if err := client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		util.Log.Error(err)
	}
}

func loginServer(c *gin.Context) {
	util.Log.Debug("aceept a cient !")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.Log.Error(err)
		return
	}
	client := &Client{
		conn:  conn,
		send:  make(chan []byte, 256),
		Games: make(map[string]struct{}),
	}
	go client.writePump()
	client.readPump()
}
