package game

import (
	//. "go/cmkj_server_go/game"
	"go/cmkj_server_go/models"
	//"go/cmkj_server_go/network"
	"go/cmkj_server_go/util"

	//"strconv"
	"strings"
	//"sync"
	//"time"
)

//DoManageCmd  管理命令
func DoManageCmd(game IGame, c *Client, cmd, tableid int, args, arg string) byte {
	util.Log.Debugf("CMD = %d", cmd)
	table := game.FindTable(tableid)
	if cmd == 1000 {
		return ManageIn(game, c, args, arg, tableid)
	}
	if table.Conf.Manager != c.P {
		if cmd == 999 {
			return MoniIn(game, c, args, arg, tableid)
		}
		if cmd == 1011 {
			table.Statu.WinLose = 0
			game.SendStatusToAll(0, tableid)
		}
		return 0
	}
	switch cmd {
	case 1001:
		return Start(game, tableid, args, arg)
	case 1002:
		game.Stop(tableid)
	case 1003: //扑克
		if table.Statu.Status == util.S_START {
			if strings.EqualFold(arg, "") {
				table.Statu.Time = 0
				table.Statu.Status = util.S_STOP
			}
		}
		table.Statu.Poker = args
		game.SendStatusToAll(0, tableid)
		util.Log.Debug("poker")
	case 1004: //结果
		game.DoResult(tableid, args, arg)
	case 1005: //下一场
		table.Statu.Status = util.S_SHUFFLE
		table.Statu.Shoe++
		table.Statu.Game = 0
		table.His.Init()
		for i := range table.Statu.Bet {
			table.Statu.Bet[i] = 0
		}
		for i := range table.Statu.VirtualBet {
			table.Statu.VirtualBet[i] = 0
		}
		game.SendStatusToAll(1, tableid)
		table.DealerServe.Shuffle(args, table.Statu.Shoe)
		models.UpdateStageHistoryStage(table.Game, tableid, table.Statu.Shoe, table.Statu.Game,
			table.Statu.Status, "", table.His.Counts)
	case 1006: //清零
		table.Statu.Status = util.S_SHUFFLE
		table.Statu.Shoe = 1
		table.Statu.Game = 0
		table.His.Init()
		game.SendStatusToAll(1, tableid)
		table.DealerServe.Shuffle(args, table.Statu.Shoe)
		models.UpdateStageHistoryStage(table.Game, tableid, table.Statu.Shoe, table.Statu.Game,
			table.Statu.Status, "", table.His.Counts)
	case 1007: //修改结果
		return game.DoChangeResult(tableid, args, arg)
	case 1019: //补充结果
		return game.DoSyncResult(tableid, args, arg)
	default:
		return 1
	}
	return 0
}

//Start 牌局开始
func Start(game IGame, tableid int, args, arg string) byte {
	game.CleanOrder(tableid)
	return game.Start(tableid, args, arg)
}

//DoManageLogin 管理员登录(荷官或者监控人员)
func DoManageLogin(game IGame, c *Client, cmd, tableid int, args, arg string) byte {
	switch cmd {
	case 1000:
		return ManageIn(game, c, args, arg, tableid)
	case 999:
		return MoniIn(game, c, args, arg, tableid)
	default:
		return 1
	}

}

//ManageIn 荷官登录
func ManageIn(game IGame, c *Client, pwd, dealer string, tableid int) byte {
	util.Log.Debug(pwd)
	table := game.FindTable(tableid)
	pwd = strings.TrimSuffix(pwd, "#")
	pwds := strings.Split(pwd, "#")
	if len(pwds) < 1 {
		return 1
	}
	if !strings.EqualFold(pwds[0], table.Conf.Password) {
		c.P = nil
		return 1
	}
	if c.P == nil {
		c.P = NewPlayer()
	}
	c.P.Home = table.Num
	c.P.Type = -1
	c.P.Send = c.Send
	if table.Conf.Manager != nil {
		if table.Conf.Manager.Send != nil && table.Conf.Manager.Send != c.P.Send {
			table.Conf.Manager.Send <- nil
		}
	}
	table.Statu.IsOpen = 1
	table.Conf.Manager = c.P
	table.Conf.Dealer = dealer
	util.Log.Debug(dealer)
	strs := strings.Split(dealer, "#")
	if len(strs) > 6 {
		//历史记录同步
		game.SyncHistory(strs, tableid)
	} else {
		//发送大厅消息
		game.SendHallMsg(table.MakeMsg(false))
	}
	game.EnterRoom(c.P, -1, tableid)
	//更新stage历史记录
	models.UpdateStageHistoryDealer(table.Game, table.Num, true, table.Conf.Dealer)
	if len(pwds) > 1 {
		//dealerService 处理
		table.DealerServe.MannageIn(table.Conf.Dealer, pwds[1])
	}
	return 0
}

//MoniIn 监控人员登录
func MoniIn(game IGame, c *Client, pwd, dealer string, tableid int) byte {
	table := game.FindTable(tableid)
	if !strings.EqualFold(table.Conf.Password, pwd) {
		c.P = nil
		return 1
	}
	if c.P == nil {
		c.P = NewPlayer()
	}
	c.P.Home = table.Num
	c.P.Type = -1
	c.P.Send = c.Send

	game.EnterRoom(c.P, c.P.Type, tableid)
	game.AddMoni(c.P, tableid)
	//向监控者发送本桌玩家信息
	go SendPlayerMsg(c.P, table)
	return 0
}
