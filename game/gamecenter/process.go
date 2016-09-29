package main

import (
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"

	"fmt"
	"strconv"
	"strings"
)

func readRimGitAndChips(id int, chips string, roulette bool) (string, byte) {
	if false == roulette {
		rs, re := models.ReadVideoChip(id)
		if re != 0 {
			return "", re
		}
		return strconv.Itoa(rs.LowerLimit) + ":" + strconv.Itoa(rs.UpperLimit) + ":" + chips + ":" + rs.ChipOptional, re
	}
	rs, re := models.ReadRouletteChip(id)
	if re != 0 {
		return "", re
	}
	return rs.LowerLimit + ":" + rs.UpperLimit + ":" + chips + ":" + rs.ChipOptional, re

}

func login(p *Player, pwd string) byte {
	//读个人信息
	member, re := models.ReadMember(p.UserName, pwd)
	if re != 0 {
		return re
	}
	p.UserID = member.UserID
	p.Blance = member.CurMoney1
	p.Limits = member.Limits
	p.MoneySort = member.MoneySort
	p.PreID = member.PreID
	p.Nickname = member.TrueName
	//读维护信息
	util.Log.Debug(p)
	util.Log.Debug(p.PreID)
	opt := models.ReadPreOpt(p.PreID)
	if opt.WhKind == 1 {
		return 6 //维护
	}
	p.IsTip = byte(opt.XfKind)
	isChat, err := strconv.Atoi(opt.Email)
	if err == nil {
		p.IsChat = byte(isChat)
		fmt.Println("is chat")
	}
	//读筹码信息
	strs1 := strings.Split(member.VideoBetLimitIDs, ",")
	limit1, _ := strconv.Atoi(strs1[0])
	strs2 := strings.Split(member.RouletteBetLimitIDs, ",")
	limit2, _ := strconv.Atoi(strs2[0])

	var chips1, chips2 string
	chips1, re = readRimGitAndChips(limit1, member.ChipVideo, false)
	if re != 0 {
		return re
	}
	chips2, re = readRimGitAndChips(limit2, member.ChipRoulette, true)
	if re != 0 {
		return re
	}
	p.Chips = chips1 + "|" + chips2

	re = models.UpdateOlineOrExit(p.UserName, util.DO_ONLINE)
	if re != 0 {
		return re
	}
	return 0
}

func getGame(way string) (gameID int, tableID int) {
	bytes := []byte(way)
	g, _ := strconv.Atoi(string(bytes[:2]))
	t, _ := strconv.Atoi(string(bytes[2:]))
	return g, t
}

func getOutPlayer(name string) {
	c, ok := DefaultCenter.clients[name]
	if false == ok {
		return
	}
	c.getoutplayer(0)
}
