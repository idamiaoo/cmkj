package game

import (
	"go/cmkj_server_go/util"
)

//MakePlayerMsg 生成玩家个人信息消息buf
func MakePlayerMsg(p *Player, tableid int) []byte {
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_PLAYER1)
	buf.WriteUTF(p.Name)
	buf.WriteDouble(p.Balance)
	buf.WriteByte(byte(tableid))
	buf.WriteShort(int16(p.Vt))
	buf.WriteByte(byte(p.Seat))

	pb := p.Bets[tableid]
	if pb.WinLost2 != 0 {
		buf.WriteDouble(pb.WinLost2)
		pb.WinLost2 = 0
	} else {
		buf.WriteDouble(pb.WinLost)
	}
	buf.WriteDouble(pb.TotalBet)
	buf.WriteUTF(pb.Betstr)
	buf.WriteByte(byte(p.IsTips))
	return buf.Bytes()
}

//SendPlayerMsgToMoni 发送玩家信息给玩家本人和监控人员
func SendPlayerMsgToMoni(p *Player, table *Table) {
	//util.Log.Debug(p.Balance)
	msg := MakePlayerMsg(p, table.Num)
	for _, moni := range table.MoniLst {
		moni.Send <- msg
	}
	p.Send <- msg
}

//SendExitMsgToMoni 发送玩家退出消息给玩家本人和监控人员
func SendExitMsgToMoni(p *Player, table *Table) {
	buf := util.NewByteBuffer()
	buf.WriteShort(19002)
	buf.WriteUTF(p.Name)
	buf.WriteDouble(p.Balance)
	buf.WriteByte(byte(table.Num))
	buf.WriteShort(-1)
	buf.WriteByte(byte(p.Platform))
	msg := buf.Bytes()
	for _, moni := range table.MoniLst {
		moni.Send <- msg
	}

}

//SendPlayerMsg  发送某桌所有玩家信息给某个监控员
func SendPlayerMsg(moni *Player, table *Table) {
	for _, p := range table.PlayerList {
		moni.Send <- MakePlayerMsg(p, table.Num)
	}
}

//SendGetResultMsg ....
func SendGetResultMsg(p *Player, num, stage, game int) {
	buf := util.NewByteBuffer()
	buf.WriteShort(int16(10409))
	buf.WriteByte(byte(num))
	buf.WriteInt(int32(stage))
	buf.WriteInt(int32(game))

	msg := buf.Bytes()
	p.Send <- msg
}
