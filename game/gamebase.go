package game

import (
	//. "go/cmkj_server_go/game"
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/network"
	"go/cmkj_server_go/util"

	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//Base 游戏服务
type Base struct {
	Max     int                 //牌桌数量
	Tables  map[int]*Table      //牌桌
	HClient *network.HallClinet //大厅通信客户端
}

//FindTables 获取游戏所有的牌桌
func (game *Base) FindTables() map[int]*Table {
	return game.Tables
}

//FindTable 获取某个牌桌
func (game *Base) FindTable(tableID int) *Table {
	table, ok := game.Tables[tableID]
	if ok {
		return table
	}
	return nil
}

//UpdatePlayerss 更新玩家信息
func (game *Base) UpdatePlayerss(tableID int, list []*PlayerStatus, islock bool) {
	table := game.Tables[tableID]
	if table == nil {
		return
	}
	if islock {
		for _, ps := range list {
			p, ok := table.PlayerList[ps.Name]
			if ok {
				p.IsLock = ps.Islock
			}
		}
		return
	}

	for _, ps := range list {
		p, ok := table.PlayerList[ps.Name]
		if ok {
			if !util.DecimalEqual(p.Balance, ps.Balance) {
				p.Balance = ps.Balance
				p.Bets[tableID].IsBet = true
			}
		}
	}
}

//Start 游戏开始
func (game *Base) Start(tableID int, args, arg string) byte {
	table := game.Tables[tableID]
	args = strings.TrimSuffix(args, "-")
	ss := strings.Split(args, "-")
	var stage, ga int
	haveway := false
	if len(ss) > 1 {
		stage, _ = strconv.Atoi(ss[0])
		ga, _ = strconv.Atoi(ss[1])
		haveway = true
	}
	var Time int
	if len(ss) > 2 {
		Time, _ = strconv.Atoi(ss[2])
	}
	if Time > 0 {
		atomic.StoreInt32(&table.Statu.Time, int32(Time))
		atomic.StoreInt32(&table.Conf.Time, int32(Time))
	} else {
		atomic.StoreInt32(&table.Statu.Time, table.Conf.Time)
	}
	if atomic.LoadInt32(&table.Statu.Status) == util.S_START {
		return 0
	}
	if atomic.LoadInt32(&table.Statu.Status) == util.S_STOP {
		atomic.StoreInt32(&table.Statu.Status, int32(util.S_START))
		return 0
	}
	if haveway && table.Statu.Shoe != stage && table.Statu.Game != ga {
		util.Log.Debug("haveway")
		table.Statu.Shoe = stage
		table.Statu.Game = ga
		table.His.Ways = arg
		table.His.Last = "q" + arg
		table.His.Counts = waysCountWay(table.His.Tj, arg)
		models.UpdateStageHistoryStage(table.Game, tableID, table.Statu.Shoe,
			table.Statu.Game, table.Statu.Status, table.His.Ways, table.His.Counts)
		game.SendStatusToAll(1, tableID)
	}
	atomic.StoreInt32(&table.Statu.Status, int32(util.S_START))
	table.Statu.Game++
	table.Statu.StartTime = time.Now()
	atomic.StoreInt32(&table.Statu.PayTime, 0)
	table.Statu.GameID = int64(table.Game+table.Num)*1000000000000000 + table.Statu.StartTime.UnixNano()/1000000
	models.UpdateStatusID(table.Game, tableID, table.Statu.Shoe, table.Statu.Game,
		table.Statu.Status, table.Statu.Time, table.Statu.GameID)
	//发送大厅消息
	util.Log.Debug(table.Statu.Time)
	game.SendHallMsg(table.MakeMsg(false))
	//clearOrder
	//game.CleanOrder(tableID)
	table.IsWriteDb = false

	return 0
}

//Stop 下注停止
func (game *Base) Stop(tableID int) {
	table := game.Tables[tableID]
	if atomic.LoadInt32(&table.Statu.Status) != util.S_START {
		return
	}
	atomic.StoreInt32(&table.Statu.Status, int32(util.S_STOP))
	if atomic.LoadInt32(&table.Statu.Time) > 0 {
		atomic.StoreInt32(&table.Statu.Time, 0)
		game.SendStatusToAll(0, tableID)
	}
}

//SendGroupMsg 给牌桌所有发玩家发送消息
func (game *Base) SendGroupMsg(msg []byte, tableID int) {
	table, ok := game.Tables[tableID]
	if false == ok {
		return
	}
	for _, cli := range table.PlayerList {
		cli.Send <- msg
	}
	if table.Conf.Manager != nil {
		table.Conf.Manager.Send <- msg
	}
}

//SendStatus 给单个玩家发送牌桌信息
func (game *Base) SendStatus(p *Player, tableID int) {
	util.Log.Debug("send status:", p.Send)
	p.Send <- game.Tables[tableID].Conf.MakeMsg()
	util.Log.Debug("send status")
	p.Send <- game.Tables[tableID].His.MakeMsg(2)
	p.Send <- game.Tables[tableID].Statu.MakeMsg()
}

//SendStatusToAll 给所有玩家发送牌桌信息
func (game *Base) SendStatusToAll(way, tableID int) {
	table := game.Tables[tableID]
	game.SendGroupMsg(table.Statu.MakeMsg(), tableID)
	if way > 0 {
		game.SendGroupMsg(table.His.MakeMsg(way), tableID)
		//发送大厅消息
		game.SendHallMsg(table.MakeMsg(true))
	}
}

//SendPlayerMsg 给所有玩家发送牌桌信息
func (game *Base) SendPlayerMsg(tableID int) {
	table := game.Tables[tableID]
	for _, p := range table.PlayerList {
		if p.Bets[tableID].IsBet {
			p.Bets[tableID].IsBet = false
			SendPlayerMsgToMoni(p, table)
		}
	}
}

//SendHallMsg 发送游戏大厅消息
func (game *Base) SendHallMsg(msg []byte) {
	game.HClient.SendMsg(msg)
}

//AddClient 牌桌增加玩家
func (game *Base) AddClient(p *Player, tableID int) {
	table := game.Tables[tableID]
	table.PlayerMutex.Lock()
	defer table.PlayerMutex.Unlock()
	_, ok := table.PlayerList[p.Name]
	if ok {
		//
		util.Log.Error("cli exists")
	}
	table.PlayerList[p.Name] = p
}

//RemoveClient 牌桌移除玩家
func (game *Base) RemoveClient(p *Player, tableID int) {
	table := game.Tables[tableID]
	table.PlayerMutex.Lock()
	defer table.PlayerMutex.Unlock()
	_, ok := table.PlayerList[p.Name]
	if false == ok {
		//
		util.Log.Error("cli non exists")
	}
	delete(table.PlayerList, p.Name)
}

//AddMoni 牌桌增加监控员
func (game *Base) AddMoni(p *Player, tableID int) {
	table := game.Tables[tableID]
	table.MoniMutex.Lock()
	defer table.MoniMutex.Unlock()
	_, ok := table.MoniLst[p.Name]
	if ok {
		//
		util.Log.Error("cli exists")
	}
	table.MoniLst[p.Name] = p
}

//RemoveMoni 监控员离开
func (game *Base) RemoveMoni(p *Player, tableID int) {
	table := game.Tables[tableID]
	table.MoniMutex.Lock()
	defer table.MoniMutex.Unlock()
	_, ok := table.MoniLst[p.Name]
	if false == ok {
		//
		util.Log.Error("cli non exists")
	}
	delete(table.MoniLst, p.Name)
}

//EnterRoom 进入牌桌
func (game *Base) EnterRoom(p *Player, t, tableID int) {
	util.Log.Debug(*p)
	table := game.Tables[tableID]
	p.Home = table.Num
	game.SendStatus(p, tableID)
	if t != -1 {
		p2, ok := table.PlayerList[p.Name]
		if ok {
			if p2 != p {
				p2.Send <- nil
				game.RemoveClient(p, tableID)
			}
			p.Bets[tableID] = p2.Bets[tableID]
		}
		game.AddClient(p, tableID)
		table.Statu.OlineNo = len(table.PlayerList)
		//发送消息
		SendPlayerMsgToMoni(p, table)
	}
}

//LeaveRoom 离开牌桌
func (game *Base) LeaveRoom(p *Player, loginOff bool, tableID int) {
	table := game.Tables[tableID]
	game.RemoveClient(p, tableID)
	table.Statu.OlineNo = len(table.PlayerList)
	if loginOff {
		if table.Conf.Manager == p {
			table.Conf.Manager = nil
			table.Conf.Dealer = ""
			table.Statu.IsOpen = 0
			//更新数据库游戏状态
			models.UpdateStageHistoryDealer(table.Game, table.Num, false, "")
			//向大厅发送游戏状态数据
			game.SendHallMsg(table.MakeMsg(false))
		}
	}
	if p.Type == -1 {
		game.RemoveMoni(p, tableID)
	} else {
		if len(table.MoniLst) > 0 {
			//给监控者发送信息
			SendExitMsgToMoni(p, table)
		}
	}
	p.Home = 0
}

//DoBet 下注
func (game *Base) DoBet(p *Player, name string, t, counts, tableID int, body []byte) byte {

	return 0
}

//DoResult 结算
func (game *Base) DoResult(tableID int, result, poker string) byte {
	return 0
}

//DoChangeResult 改变结果
func (game *Base) DoChangeResult(tableID int, result, poker string) byte {
	return 0
}

//DoSyncResult 同步结果
func (game *Base) DoSyncResult(tableID int, result, poker string) byte {
	return 0
}

//CleanOrder 清空注单表
func (game *Base) CleanOrder(tableID int) {
	table := game.Tables[tableID]
	for _, p := range table.PlayerList {
		if p.Bets[tableID].TotalBet == 0 {
			continue
		}
		p.Bets[tableID] = NewBet()
	}

	table.Orders.Clear()

	table.Statu.ClearBet()

	table.Statu.Poker = ""
	table.Statu.Result = ""
}

//SyncHistory 同步历史记录
func (game *Base) SyncHistory(strs []string, tableID int) {
	table := game.Tables[tableID]
	table.Conf.Dealer = strs[0]
	shoe, _ := strconv.Atoi(strs[1])
	times, _ := strconv.Atoi(strs[2])
	stat, _ := strconv.Atoi(strs[3])

	table.Statu.Time = 0

	isDef := false

	if shoe != table.Statu.Shoe || times != table.Statu.Game {
		isDef = true
		if stat == util.S_STOP || stat == util.S_START {
			table.Statu.Status = util.S_STOP
			table.Statu.GameID = 0
		} else {
			table.Statu.Status = util.S_OVER
		}
	} else {
		if table.Statu.Status == util.S_START || table.Statu.Status == util.S_STOP {
			if stat == util.S_PAYOFF || stat == util.S_OVER {
				isDef = true
				table.Statu.Status = util.S_OVER
			} else {
				table.Statu.Status = util.S_STOP
			}
		} else {
			table.Statu.Status = util.S_OVER
		}
	}

	if isDef {
		table.His.Ways = strs[5]
		table.His.Last = "q" + strs[5]
		table.His.Counts = waysCountWay(table.His.Tj, strs[5])

		if table.Orders.Len() > 0 {
			table.Orders.Clear()
		}
		table.Statu.Shoe = shoe
		table.Statu.Game = times
		models.UpdateStageHistoryStage(table.Game, table.Num, table.Statu.Shoe, table.Statu.Game,
			table.Statu.Status, table.His.Ways, table.His.Counts)
		game.SendStatusToAll(1, tableID)
	} else {
		game.SendHallMsg(table.MakeMsg(false))
	}
}

//TimeRun 牌桌后台定时器
func (game *Base) TimeRun() {
	tables := game.FindTables()
	for i := range tables {
		go game.run(i)
	}
}

func (game *Base) run(tableID int) {
	for {
		<-time.After(time.Second * time.Duration(1))
		game.gogo(tableID)
	}
}

func (game *Base) gogo(tableID int) {
	table := game.FindTable(tableID)
	if atomic.LoadInt32(&table.Statu.Time) != 0 {

		if atomic.AddInt32(&table.Statu.Time, -1) <= 0 {
			atomic.StoreInt32(&table.Statu.Status, int32(util.S_STOP))
		}
		game.SendStatusToAll(0, tableID)
	} else if atomic.LoadInt32(&table.Statu.PayTime) != 0 {
		if atomic.AddInt32(&table.Statu.PayTime, -1) <= 0 {
			atomic.StoreInt32(&table.Statu.Status, int32(util.S_OVER))
			models.UpdateStatusTime(table.Game, tableID, table.Statu.Status, table.Statu.Time)
			game.CleanOrder(tableID)
			game.SendStatusToAll(0, tableID)
		}
	}
}

//VirtualBetRun 虚拟投注
func (game *Base) VirtualBetRun() {
	tables := game.FindTables()
	for tabledID := range tables {
		go game.doVirtualBet(tabledID)
	}
}

func (game *Base) doVirtualBet(tableID int) {
	table := game.FindTable(tableID)
	virtualStr := strings.TrimSuffix(table.Conf.VirtualStr, "#")
	d, err := util.SplitInt(virtualStr, "#")
	if err != nil {
		util.Log.Fatal(err)
	}
	if len(d) < 1 {
		util.Log.Fatalf("Virtual err len(d)=%d", len(d))
	}
	for {
		<-time.After(time.Duration(1) * time.Second)
		if table.Statu.Time > 3 {
			if table.Conf.VirtualBet == 1 {
				source := rand.NewSource(time.Now().Unix())
				random := rand.New(source)
				sleep := random.Intn(d[1]-d[0]) + d[0]
				<-time.After(time.Duration(sleep) * time.Millisecond)
				for i := 0; i < 5; i++ {
					index := i*2 + 2
					if random.Intn(2) == 1 {
						sleep = random.Intn(d[index+1]-d[index]) + d[index]
						table.Statu.VirtualBet[i] += (sleep / 10 * 10)
					}
				}
			}
		}
	}
}
