package main

import (
	"go/cmkj_server_go/conf"
	Game "go/cmkj_server_go/game"
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/network"
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//BjlServer 百家乐游戏服务端
type BjlServer struct {
	Game.Base
	VTables map[int]*vipRooms //虚拟5人桌
}

//Bjl 全局游戏服务实例
var Bjl *BjlServer

//NewBjlTable 新建百家乐牌桌
func NewBjlTable(num int, config string) *Game.Table {
	table := &Game.Table{
		Statu:      Game.NewStatus(6),
		Conf:       Game.NewConfig(GBetmaxid),
		His:        Game.NewHistory(7),
		Game:       GGame,
		Num:        num,
		PlayerList: make(map[string]*Game.Player),
		MoniLst:    make(map[string]*Game.Player),
	}

	table.Statu.Num = num
	table.Conf.Num = num
	table.His.Num = num
	if !table.Conf.LoadConfig(config) {
		util.Log.Fatal("load table config failed")
	}
	if !table.GetStageHistory() {
		util.Log.Fatal("get stage history failed")
	}
	if !table.ReadGamePoker() {
		util.Log.Fatalf("table[%d] read game poker failed", table.Num)
	}

	table.Orders = Game.LoadTmpOrders(table.Game, 0, table.Num, table.Statu.Shoe, table.Statu.Game)
	util.Log.Info(table.Orders.Len())
	util.Log.Info(table.Num)
	table.DealerServe = Game.NewDealerService(table.Game, table.Num, true)
	return table
}

//LoadTables 根据配置创建牌桌
func LoadTables() map[int]*Game.Table {
	tables := make(map[int]*Game.Table)
	config := conf.Conf.String("rooms")
	configs := strings.Split(config, ",")
	if len(configs) <= 0 {
		util.Log.Fatal("no table config")
	}
	for i, name := range configs {
		tables[i+1] = NewBjlTable(i+1, name)
	}
	return tables
}

//NewBjl 新增百家乐实例
func NewBjl() *BjlServer {
	bjl := &BjlServer{
		Base: Game.Base{
			Tables:  LoadTables(),
			HClient: network.NewhallClient(),
		},
		VTables: make(map[int]*vipRooms),
	}
	for k := range bjl.Tables {
		bjl.VTables[k] = newVipRooms()
	}
	bjl.Max = len(bjl.Tables)
	util.Log.Info(bjl.Tables[1].Conf.Rate)
	bjl.HClient.Start()
	bjl.TimeRun()
	bjl.VirtualBetRun()
	return bjl
}

//EnterRoom 进入牌桌
func (game *BjlServer) EnterRoom(p *Game.Player, t, tableID int) {
	util.Log.Debug("enter room :", t, tableID)
	table := game.Tables[tableID]
	p.Home = table.Num
	game.SendStatus(p, tableID)
	if t != -1 {
		util.Log.Debug("enter room")
		p2, ok := table.PlayerList[p.Name]
		if ok {
			if p2 != p {
				p2.Send <- nil
				game.RemoveClient(p2, tableID)
			}
			p.Bets[tableID] = p2.Bets[tableID]
			//退出原vip房间
			game.VTables[tableID].leaveRoom(p2)

		}
		//获取历史下注信息
		game.GetHistoryBet(p, tableID)
		game.AddClient(p, tableID)
		table.Statu.OlineNo = len(table.PlayerList)
		//发送消息
		Game.SendPlayerMsgToMoni(p, table)
		if t == 1 {
			//进入vip房间
			game.VTables[tableID].inRoom(p)
			//vip房间下注处理
			game.VTables[tableID].bet(p, true)
			util.Log.Debug("enter room")
		}
	}
	util.Log.Debug("enter room")
}

//LeaveRoom 离开牌桌
func (game *BjlServer) LeaveRoom(p *Game.Player, loginOff bool, tableID int) {
	game.VTables[tableID].leaveRoom(p)
	//game.CleanOrder()
	game.Base.LeaveRoom(p, loginOff, tableID) //调用
}

//CleanOrder 清理注单
func (game *BjlServer) CleanOrder(tableID int) {
	game.Base.CleanOrder(tableID)
	game.VTables[tableID].clearbet()

}

//GetHistoryBet 获取下注记录
func (game *BjlServer) GetHistoryBet(p *Game.Player, tableID int) {
	pb := p.Bets[tableID]
	if pb.BetData != nil {
		for i := 0; i < len(pb.BetData); i++ {
			pb.BetData[i] = 0
		}
	}
	pb.TotalBet = 0
	for _, bo := range game.Tables[tableID].Orders.Orders() {
		if strings.EqualFold(bo.UserName, p.Name) {
			pb.BetData = joinBetWithSlice(pb.BetData, bo.BettingAmount,
				bo.GameBettingKind-GBetindex)
			pb.TotalBet = bo.BettingAmount
			bo.P = p
		}
	}
	if pb.TotalBet > 0 {
		pb.Betstr = converBetToStr(pb.BetData)
	}

}

//settlement 结算
func (game *BjlServer) settlement(tableID int, orders []*Game.BettingOrder, results []byte, way string, newid, oldid int64, issync bool) {
	util.Log.Debug(len(orders))
	table := game.Tables[tableID]
	pay := NewPay()
	rate := table.Conf.Rate
	for _, b := range orders {
		util.Log.Debug(b)
		pay.Init()
		index := b.GameBettingKind - GBetindex - 1
		util.Log.Debug(index)
		if pay.getPayBool(results, b.BettingAmount, index, b.GameKind, rate) {
			b.CompensateRate = rate[index]
			util.Log.Debug("true")
		}
		money := pay.Pay
		models.MoneyAdd(b.UserName, &money, 1, b.GameKind, b.OrderNumber, b.UserID, b.PlatformID)
		b.Balance = money
		b.WinLoseAmount = pay.WinLose
		if pay.WinLose < 0 {
			b.ResultType = 1
		} else if pay.WinLose == 0 {
			b.ResultType = 3
		} else {
			b.ResultType = 2
		}
		b.GameResult = way
		b.GameResultContent = pay.Pay
		b.GameRecordID = newid
		b.GameBettingContent = strconv.Itoa(b.GameBettingKind) + "^" + util.DecimalFormat(b.BettingAmount) +
			"^" + util.DecimalFormat(b.WinLoseAmount) + "^"
		table.Statu.WinLose -= b.WinLoseAmount
		table.DealerServe.AddWinLose(1, b.BettingAmount, b.WinLoseAmount, 0)
		if b.P != nil && b.P.Online {
			b.P.Balance = b.Balance
			pb := b.P.Bets[tableID]
			pb.Betstr = ""
			pb.Pay += pay.Pay
			pb.WinLost += b.WinLoseAmount
			pb.IsBet = true
			game.VTables[tableID].bet(b.P, false)
		}
	}
	recordrun := &Game.RecordRun{
		Game:    game,
		TableID: tableID,
	}
	if !issync {
		go recordrun.Run()
	}
	//结算记录
	Game.ToOrderAll(orders)
	//删除临时表
	models.DeleteTmpWithRid(oldid)
}

//DoBet 下注
func (game *BjlServer) DoBet(p *Game.Player, name string, Type, counts, tableID int, body []byte) byte {
	util.Log.Debug("dobet:", counts, Type, p, name)
	table := game.Tables[tableID]
	msg := util.NewByteBufferWith(body)
	if atomic.LoadInt32(&table.Statu.Status) != util.S_START {
		return 2
	}
	if int64(table.Conf.Time*1000)+table.Statu.StartTime.Unix() < time.Now().Unix() {
		return 2
	}
	/*
		if p.Bets[tableID].Last+200 > time.Now().UnixNano()/1000000 { //haomiao
			return 1
		}
		p.Bets[tableID].Last = time.Now().UnixNano() / 1000000
	*/
	switch Type {
	case GBjl:
	case GMianyong:
	case GLianhuan:
	case 115:
	default:
		return 1
	}
	if counts < 1 || counts > 256 {
		return 1
	}

	if p == nil || !strings.EqualFold(name, p.Name) {
		return 3
	} else if p.IsLock == 1 {
		return 5
	} else if p.Maxprofit > 0 && p.Balance > p.Maxprofit {
		return 6
	}
	var (
		index int
		total float64
	)
	indexs := make([]int, counts)
	totals := make([]float64, counts)

	pb := p.Bets[tableID]
	util.Log.Debug(pb)
	for i := 0; i < counts; i++ {
		util.Log.Debug("for")
		indexs[i] = int(msg.ReadShort())
		totals[i] = msg.ReadDouble()
		util.Log.Debug(indexs[i], totals[i])
		if totals[i] < 0 {
			return 1
		}
		total += totals[i]
		index = indexs[i] - GBetindex
		//util.Log.Debug(index)
		if index < 1 || index > GBetmaxid {
			return 1
		}
		if index > 5 && index < 26 {
			return 1
		}
		index--
		var t float64
		if pb.BetData != nil {
			t = pb.BetData[index]
		}
		if t+totals[i] > float64(p.Limit)/table.Conf.Rate[index] {
			return 6
		}

	}
	if total <= 0 {
		return 1
	}
	var err byte
	addTime := time.Now()
	util.Log.Debug(indexs, totals)
	for i := 0; i < counts; i++ {
		index = indexs[i]
		total = totals[i]

		id := util.UIDGenerator.GetID(GThisnaem)
		sid := strconv.FormatInt(id, 10)
		//数据库扣费操作
		money := -total
		ok := models.MoneyAdd(p.Name, &money, 0, Type, sid, p.ID, p.Platform)
		if ok == 1 {
			err = 4
			break
		}
		p.Balance = money
		pb.TotalBet += total

		bo := Game.NewBettingOrder()
		bo.ID = id
		bo.TableID = tableID
		bo.OrderNumber = sid
		bo.GameKind = Type
		bo.GameBettingKind = index
		bo.BettingAmount = total
		bo.GameRecordID = table.Statu.GameID
		bo.GameNameID = table.Game
		bo.P = p
		bo.UserID = p.ID
		bo.UserName = p.Name
		bo.Balance = p.Balance
		bo.PlatformID = p.Platform
		bo.AddTime = addTime
		bo.IP = p.IP
		bo.Path = p.PreSequence
		bo.StartTime = table.Statu.StartTime
		bo.GameBettingContent = strconv.Itoa(index) + "^" + util.DecimalFormat(total) + "^" + "0" + "^" + ","
		bo.GameResult = ""
		bo.Stage = table.Statu.Shoe
		bo.Inning = table.Statu.Game
		tmp := Game.OrderConvertToTmp(bo)
		if !tmp.Create() {
			//资金回退
			money = total
			models.MoneyAdd(p.Name, &money, -1, Type, sid, p.ID, p.Platform)
			return 7
		}
		table.Orders.Add(bo)
		index -= GBetindex
		pb.BetData = joinBetWithSlice(pb.BetData, total, index)

		if index == GBetindex {
			index = 6
		}
		if index < 7 {
			table.Statu.AddBet(index-1, int(total))
		}
	}
	pb.Betstr = converBetToStr(pb.BetData)
	game.VTables[tableID].bet(p, true)
	Game.SendPlayerMsgToMoni(p, table)
	return err
}

//DoResult 结算命令处理
func (game *BjlServer) DoResult(tableID int, result, poker string) byte {
	table := game.Tables[tableID]
	pokerTmp := strings.TrimSuffix(poker, "-")
	str := strings.Split(pokerTmp, "-")
	util.Log.Debug(str)
	util.Log.Debug(len(str))
	if len(str) < 3 {
		return 1
	}
	c1, _ := strconv.Atoi(str[0])
	c2, _ := strconv.Atoi(str[1])
	poker = str[2]
	if c1 != table.Statu.Shoe || c2 != table.Statu.Game {
		return 1
	}
	results := converResult(result)
	if len(results) < 6 {
		return 1
	}
	results = getNewResult(poker)
	if results == nil {
		return 1
	}
	if atomic.LoadInt32(&table.Statu.Status) != util.S_STOP {
		return 2
	}

	//t := time.Now().Unix()

	atomic.StoreInt32(&table.Statu.Status, int32(util.S_PAYOFF))
	table.Statu.Poker = poker

	way := converWay(results)
	table.His.Last = way
	table.His.Ways += table.His.Last
	table.His.Counts = waysCount(table.His.Tj, results)
	table.His.Pokers = table.His.Pokers + poker + ","
	table.His.Poker = poker
	table.Statu.Result = way

	//把结果记录至数据库
	newid := models.WriteResult(table.Game, table.Num, table.Statu.Shoe, table.Statu.Game,
		table.Statu.StartTime, way, poker, 0, 0)
	if newid <= 0 {
		newid = models.WriteResult(table.Game, table.Num, table.Statu.Shoe, table.Statu.Game,
			table.Statu.StartTime, way, poker, 0, 0)
	}
	table.Statu.PreID = newid
	util.Log.Debug("preid:", newid)

	//异步结算
	go game.settlement(tableID, table.Orders.Orders(), results, way, newid, table.Statu.GameID, false)
	table.Statu.WinLose = util.Round(table.Statu.WinLose, 2)
	table.Statu.Status = util.S_PAYOFF
	table.Statu.PayTime = 6
	//记录结果
	models.UpdateStageHistoryWinLose(table.Game, tableID, table.Statu.Status, table.His.Ways, table.His.Counts, table.Statu.WinLose)

	game.SendStatusToAll(1, tableID)
	//删除临时表
	table.DealerServe.ToDealerStatus()
	table.DealerServe.ToShuffleStatus()
	return 0
}

//DoChangeResult 更改游戏结果命令
func (game *BjlServer) DoChangeResult(tableID int, result, poker string) byte {
	table := game.Tables[tableID]
	pokerTmp := strings.TrimSuffix(poker, "-")
	ss := strings.Split(pokerTmp, "-")
	var (
		stage, inning, isOld = 0, 0, 0
	)
	if len(ss) < 3 {
		return 1
	}
	stage, _ = strconv.Atoi(ss[0])
	inning, _ = strconv.Atoi(ss[1])
	poker = ss[2]

	if len(ss) > 3 {
		isOld, _ = strconv.Atoi(ss[3])
	}
	results := converResult(result)
	if len(results) < 6 {
		return 1
	}
	table.Statu.Poker = poker
	results = getNewResult(poker)
	if results == nil {
		return 1
	}
	way := converWay(results)

	if isOld == 0 && len(table.His.Ways) > 0 {
		table.His.Last = way
		table.His.Ways = table.His.Ways[:len(table.His.Ways)-1] + table.His.Last
		table.His.Last = "q" + table.His.Ways
		table.His.Counts = waysCountWay(table.His.Tj, table.His.Ways)
		flag := strings.LastIndex(",", table.His.Pokers[:len(table.His.Pokers)-1])
		table.His.Pokers = table.His.Pokers[:flag+1] + poker + ","
		table.His.Poker = table.His.Pokers
	}
	go func() {
		pay := NewPay()
		newid := table.Statu.PreID
		rate := table.Conf.Rate
		orders := Game.LoadTmpOrders(table.Game, 1, tableID, stage, inning)
		//var newid int64 = 0
		for _, b := range orders.Orders() {
			pay.Init()
			newid = b.GameRecordID
			index := b.GameBettingKind - GBetindex - 1

			if pay.getPayBool(results, b.BettingAmount, index, b.GameKind, rate) {
				b.CompensateRate = rate[index]
			}
			back := pay.WinLose - b.WinLoseAmount
			oldwin := -(b.WinLoseAmount - pay.WinLose)
			money := back
			models.MoneyAdd(b.UserName, &money, 2, b.GameKind, b.OrderNumber, b.UserID, b.PlatformID)
			b.Balance = money
			b.WinLoseAmount = pay.WinLose
			if pay.WinLose < 0 {
				b.ResultType = 1
			} else if pay.WinLose == 0 {
				b.ResultType = 3
			} else {
				b.ResultType = 2
			}
			b.GameResult = way
			b.GameResultContent = pay.Pay
			b.GameBettingContent = strconv.Itoa(b.GameBettingKind) + "^" + util.DecimalFormat(b.BettingAmount) +
				"^" + util.DecimalFormat(b.WinLoseAmount) + "^,"
			table.Statu.WinLose -= oldwin
			table.DealerServe.AddWinLose(0, 0, oldwin, 0)
			ps, ok := table.PlayerList[b.UserName]
			if ok {
				b.P = ps
				b.P.Balance = b.Balance
				pb := b.P.Bets[tableID]
				pb.WinLost2 += b.WinLoseAmount
				pb.IsBet = true
				game.VTables[tableID].bet(b.P, true)
			}
		}
		models.WriteResult(table.Game, tableID, stage, inning, table.Statu.StartTime, way, poker, newid, 1)
		game.SendPlayerMsg(tableID)
	}()
	models.UpdateStageHistoryWinLose(table.Game, tableID, table.Statu.Status, table.His.Ways, table.His.Counts, table.Statu.WinLose)
	table.Statu.Poker = ""
	if isOld == 0 {
		game.SendStatusToAll(1, tableID)
	} else {
		game.SendStatusToAll(0, tableID)
	}
	return 0
}

//DoSyncResult 未结算时重结算
func (game *BjlServer) DoSyncResult(tableID int, result, poker string) byte {
	table := game.Tables[tableID]
	pokerTmp := strings.TrimSuffix(poker, "-")
	str := strings.Split(pokerTmp, "-")
	var (
		stage, inning int = 0, 0
	)
	if len(str) < 3 {
		return 1
	}
	stage, _ = strconv.Atoi(str[0])
	inning, _ = strconv.Atoi(str[1])
	poker = str[2]
	results := converResult(result)
	if len(results) < 6 {
		return 1
	}
	results = getNewResult(poker)
	if results == nil {
		return 1
	}
	way := converWay(results)
	dd := time.Now()
	newid := models.WriteResult(table.Game, table.Num, table.Statu.Shoe, table.Statu.Game,
		dd, way, poker, 0, 2)
	//获取临时表
	orders := Game.LoadTmpOrders(table.Game, 0, tableID, stage, inning)
	if orders.Len() <= 0 {
		return 0
	}
	borders := orders.Orders()
	oldid := borders[0].GameRecordID
	go game.settlement(tableID, borders, results, way, newid, oldid, true)
	table.DealerServe.ToDealerStatus()
	table.DealerServe.ToShuffleStatus()
	return 0
}
