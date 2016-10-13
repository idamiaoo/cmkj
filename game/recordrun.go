package game

import (
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"

	"strconv"
)

//PlayerStatus 玩家状态
type PlayerStatus struct {
	Name    string
	Islock  int
	Balance float64
}

//UToPlayStatus 数据库记录转成玩家状态
func UToPlayStatus(m *models.MemberUpload) *PlayerStatus {
	islock, err := strconv.Atoi(m.Content)
	if err != nil {
		util.Log.Error(err)
		return nil
	}
	return &PlayerStatus{
		Name:   m.UserName,
		Islock: islock,
	}
}

//ToPlayStatus 数据库记录转成玩家状态
func ToPlayStatus(m *models.Member) *PlayerStatus {
	return &PlayerStatus{
		Balance: m.CurMoney1,
		Name:    m.UserName,
	}
}

//RecordRun ...
type RecordRun struct {
	Game    IGame
	TableID int
}

//NewRecordRun ...
func NewRecordRun(game IGame, tableid int) *RecordRun {
	return &RecordRun{
		Game:    game,
		TableID: tableid,
	}
}

//Run ...
func (r *RecordRun) Run() {
	r.updatePlayers()
	r.Game.SendPlayerMsg(r.TableID)
	r.checkNotPayoff()
}

func (r *RecordRun) updatePlayers() {
	util.Log.Debug("updatePlayers")
	table := r.Game.FindTable(r.TableID)
	if table == nil {
		return
	}
	members := models.GetMemberStatus(0, table.Game)
	var list []*PlayerStatus
	for _, member := range members {
		ps := UToPlayStatus(&member)
		if ps != nil {
			list = append(list, ps)
		}
	}
	r.Game.UpdatePlayerss(r.TableID, list, true)
	var name []string
	for _, p := range table.PlayerList {
		if !p.Bets[r.TableID].IsBet {
			name = append(name, p.Name)
		}
		if len(name) == 100 {
			nameClone := make([]string, len(name))
			copy(nameClone, name)
			go r.updateMoneyRun(nameClone)
			name = make([]string, 0)
		}
	}
	if len(name) > 0 {
		go r.updateMoneyRun(name)
	}

}

//查询是否有无结算的场次
func (r *RecordRun) checkNotPayoff() {
	util.Log.Debug("checkNotPayoff")
	table := r.Game.FindTable(r.TableID)
	if table.Statu.Status == 1 {

	}
	tmp := models.GetTempNotPayoff(table.Game, table.Num, table.Statu.GameID)
	if tmp == nil {
		return
	}
	var list []*BettingOrder
	for _, t := range tmp {
		l := &BettingOrder{
			Stage:  t.Stage,
			Inning: t.Inning,
		}
		list = append(list, l)
	}
	for _, b := range list {
		SendGetResultMsg(table.Conf.Manager, table.Num, b.Stage, b.Inning)
	}
}

func (r *RecordRun) updateMoneyRun(name []string) {
	util.Log.Debug("updateMoneyRun")
	//数据库操作
	members := models.UpdateMoneyIn(name)
	if members == nil {
		return
	}
	var list []*PlayerStatus
	for _, m := range members {
		list = append(list, ToPlayStatus(&m))
	}
	r.Game.UpdatePlayerss(r.TableID, list, false)
	r.Game.SendPlayerMsg(r.TableID)
}
