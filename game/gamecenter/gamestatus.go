package main

import (
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/timer"
	"go/cmkj_server_go/util"
)

type Status struct {
	GameType int
	TableID  int
	Stage    int    // 场
	Game     int    // 次
	Status   int    // 状态
	Time     int    // 倒计时
	Startime string // 开始时间
	Gameid   int64  // 局编号
	Result   string // 结果
	Poker    string // 扑克
	Paytime  int    // 结算延时
	Ways     string
	Last     string
	Counts   string //统计
	IsOpen   int
	Dealer   string //荷官
	Limit    int    //限红
	OnlineNo int    //在线人数
}

func getGamestatus(status *Status) bool {
	gamestatus := models.ReadStageHistory(status.GameType, status.TableID)
	if gamestatus == nil {
		util.Log.Debugf("table [%d] not exists", status.TableID)
		return false
	}
	status.Ways = gamestatus.Results
	status.Counts = gamestatus.Counts
	status.Dealer = gamestatus.Dealer
	status.Stage = gamestatus.Stage
	status.Game = gamestatus.Inning
	status.Status = int(gamestatus.Status)
	if gamestatus.IsOpen {
		status.IsOpen = 1
	}
	status.Limit = gamestatus.Limit
	t, ok := timer.GameTimer[int(status.GameType*10+status.TableID)]
	if true == ok {
		status.Time = t
	}
	return true
}

func (status *Status) makeMsg() []byte {
	//til.Log.Debug(status.GameType, status.TableID, status.Ways, status.Counts)
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_STATUS1)
	buf.WriteByte(byte(status.GameType))
	buf.WriteByte(byte(status.TableID))
	buf.WriteInt(int32(status.Stage))
	buf.WriteInt(int32(status.Game))
	buf.WriteByte(byte(status.Status))
	buf.WriteShort(int16(status.Time))
	buf.WriteUTF(status.Ways)
	buf.WriteUTF(status.Counts)
	buf.WriteByte(byte(status.IsOpen))
	buf.WriteUTF(status.Dealer)
	buf.WriteInt(int32(status.Limit))
	return buf.Bytes()

}
