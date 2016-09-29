package main

import (
	Game "go/cmkj_server_go/game"
	"go/cmkj_server_go/util"

	"strings"
	"sync"
)

var rmax = 5

type vipsit struct {
	id      int64
	name    string
	balance float64
	betstr  string
	s       chan []byte
}

func newVipsit() *vipsit {
	return &vipsit{}
}

type vipRoom struct {
	counts int
	vs     []*vipsit
	xy     string
	mutex  sync.Mutex
}

func newVipRoom() *vipRoom {
	room := &vipRoom{
		vs: make([]*vipsit, 0, 5),
	}
	for i := 0; i < cap(room.vs); i++ {
		room.vs = append(room.vs, newVipsit())
	}
	return room
}

func (vr *vipRoom) sendMsg() {
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_VIPROOM)
	for i, vs := range vr.vs {
		buf.WriteByte(byte(i + 1))
		buf.WriteUTF(vs.name)
		buf.WriteUTF(vs.betstr)
		buf.WriteDouble(vs.balance)
	}
	msg := buf.Bytes()
	for _, vs := range vr.vs {
		vs.s <- msg
	}
}

type vipRooms struct {
	rooms []*vipRoom
}

func newVipRooms() *vipRooms {
	return &vipRooms{
		rooms: make([]*vipRoom, 0),
	}
}

func (rooms *vipRooms) inRoom(p *Game.Player) int {
	for i, vr := range rooms.rooms {
		if vr.counts == rmax {
			continue
		}
		if !strings.EqualFold(p.Xy, vr.xy) {
			continue
		}
		//进入vip房间失败,重新寻找空位
		if !rooms.sitDown(vr, p) {
			continue
		}
		p.Vt = i + 1
	}
	vr := newVipRoom()
	rooms.sitDown(vr, p)
	vr.xy = p.Xy
	rooms.rooms = append(rooms.rooms, vr)
	p.Vt = len(rooms.rooms)
	return p.Vt
}
func (rooms *vipRooms) sitDown(vr *vipRoom, p *Game.Player) bool {
	vr.mutex.Lock()
	defer vr.mutex.Unlock()
	if vr.counts == rmax {
		return false
	}
	for i, vs := range vr.vs {
		if vs.id == 0 {
			vs.balance = p.Balance
			vs.name = p.Name
			vs.betstr = p.Bets[p.Home].Betstr
			vs.id = p.ID
			vs.s = p.Send
			p.Seat = i + 1
		}
	}
	vr.counts++
	return true
}

func (rooms *vipRooms) leaveRoom(p *Game.Player) {
	rid := p.Vt
	sid := p.Seat
	if rid <= 0 || sid <= 0 || rid >= len(rooms.rooms) {
		return
	}
	rid--
	sid--
	vr := rooms.rooms[rid]
	vr.mutex.Lock()
	defer vr.mutex.Unlock()
	vr.vs[sid] = newVipsit()
	vr.counts--
	if vr.counts < 1 {
		vr.counts = 0
		vr.xy = ""
		return
	}
	//vip房间消息
	vr.sendMsg()
}

func (rooms *vipRooms) bet(p *Game.Player, isMsg bool) {
	if p.Vt <= 0 || p.Seat <= 0 || p.Vt >= len(rooms.rooms) {
		return
	}
	rid := p.Vt
	sid := p.Seat

	rid--
	sid--
	vr := rooms.rooms[rid]
	vr.vs[sid].betstr = p.Bets[p.Home].Betstr
	vr.vs[sid].balance = p.Balance
	if isMsg {
		//vip房间消息
		vr.sendMsg()
	}
}

func (rooms *vipRooms) clearbet() {
	for _, vr := range rooms.rooms {
		if vr.counts == 0 {
			continue
		}
		for _, vs := range vr.vs {
			vs.betstr = ""
		}
		//vip 房间消息
		vr.sendMsg()
	}
}
