package models

import (
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
	"time"
)

type Member struct {
	UserID              int64   `gorm:"column:userid"`
	UserName            string  `gorm:"column:UserName"`
	UserPass            string  `gorm:"column:userpass"`
	IsUseable           string  `gorm:"column:isuseable"`
	isOline             int     `gorm:"column:IsOnline"`
	CurMoney1           float64 `gorm:"column:curMoney1"`
	VideoBetLimitIDs    string  `gorm:"column:VideoBetLimitIDs"`
	RouletteBetLimitIDs string  `gorm:"column:RouletteBetLimitIDs"`
	ChipVideo           string  `gorm:"column:ChipVideo"`
	ChipRoulette        string  `gorm:"column:ChipRoulette"`
	Limits              string  `gorm:"column:limits"`
	MoneySort           string  `gorm:"column:moneysort"`
	PreID               int64   `gorm:"column:pre_id"`
	TrueName            string  `gorm:"column:TrueName"`
	ClassID             int     `gorm:"column:classid"`
	Maxprofit           float64 `gorm:"column:winPoint"`
	IsLock              int     `gorm:"column:islock"`
	PreSequence         string  `gorm:"column:pre_sequence"`
	SessionID           string  `gorm:"column:sessionid"`
}

func (Member) TableName() string {
	return "member"
}

func ReadMember(username, pwd string) (*Member, byte) {
	var members []Member
	if err := DBEngine.Db.Where("UserName = ?", username).Find(&members).Error; err != nil {
		util.Log.Error(err)
		return nil, 1
	}
	if len(members) != 1 {
		return nil, 2
	}
	member := members[0]
	util.Log.Debugf("%v\n", member)

	if !strings.EqualFold(member.UserPass, pwd) {
		return nil, 3
	}

	if !strings.EqualFold(member.IsUseable, "1") {
		return nil, 4
	}

	if member.isOline != 0 {
		util.Log.Info("重复登录")
	}
	return &member, 0
}

func UpdateOlineOrExit(name string, opt int) byte {
	member := &Member{
		UserName: name,
	}
	if err := DBEngine.Db.Model(member).Update("isonline ", opt).Error; err != nil {
		util.Log.Error(err)
		return 1
	}
	return 0
}

func ChangeNickName(pwd, username, nickname string) byte {
	member := &Member{
		UserName: username,
		UserPass: pwd,
	}
	if err := DBEngine.Db.Model(member).Update("TrueName", nickname).Error; err != nil {
		util.Log.Error(err)
	}
	return 0
}

func UpdateLoginInfo(ip, username string, opt int) byte {
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	member := &Member{
		UserName: username,
	}
	if err := DBEngine.Db.Model(member).Updates(map[string]interface{}{
		"isonline":   opt,
		"updatetime": nowtime,
		"AddressIP":  ip,
	}).Error; err != nil {
		util.Log.Error(err)
		return 1
	}
	return 0
}

func UpdateMoney(username string, money *float64, pm []string, platform int) byte {
	var members []Member
	if err := DBEngine.Db.Where("UserName = ?", username).Find(&members).Error; err != nil {
		util.Log.Error(err)
		return 1
	}
	if len(members) <= 0 {
		return 1
	}
	*money = members[0].CurMoney1 + *money
	pm[0] = members[0].MoneySort
	pm[1] = members[0].PreSequence
	if *money < 0 {
		return 1
	}
	if err := DBEngine.Db.Exec(`update member set curMoney1 = ? where username = ?`,
		money, username).Error; err != nil {
		util.Log.Error(err)
		return 1
	}
	return 0
}

func MoneyAdd(username string, money *float64, remark0, remark1 int, remark2 string, id int64, platform int) byte {
	moneysort := make([]string, 2)
	tz := *money
	util.Log.Debug(tz, *money)
	ok := UpdateMoney(username, money, moneysort, platform)
	if ok == 0 {
		his := &XXHistory{
			UserName:  username,
			MoneySort: moneysort[0],
			BallTime:  time.Now(),
			Remark0:   strconv.Itoa(remark0),
			Remark1:   strconv.Itoa(remark1),
			Remark2:   remark2,
			Tzmoney:   tz,
			Comm:      0,
			Total:     *money,
			Path:      moneysort[1],
		}
		his.Create()
	}
	return ok
}

func UpdateMoneyIn(name []string) []Member {
	var members []Member
	if err := DBEngine.Db.Where("username in (?)", name).Find(&members).Error; err != nil {
		util.Log.Error(err)
		return nil
	}
	return members
}
