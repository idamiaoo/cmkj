package models

import (
	"go/cmkj_server_go/util"
	//"strings"
)

type MemberUpload struct {
	ID         int64  `gorm:"column:ID"`
	UserName   string `gorm:"column:UserName"`
	Type       int    `gorm:"column:Type"`
	Content    string `gorm:"column:Content"`
	Error      int    `gorm:"column:Error"`
	AddTime    string `gorm:"column:AddTime"`
	UseTime    string `gorm:"column:UseTime"`
	Uploaded   bool   `gorm:"column:Uploaded"`
	Remark0    string `gorm:"column:remark0"`
	Remark1    string `gorm:"column:remark1"`
	Remark2    string `gorm:"column:remark2"`
	CashTradID int64  `gorm:"column:CashTradeID"`
}

func (MemberUpload) TableName() string {
	return "MemberUpload"
}

func ReadMemberUpload() []MemberUpload {
	var umember []MemberUpload
	if err := DBEngine.Db.Where("type = ?", 4).Find(&umember).Error; err != nil {
		util.Log.Debug(err)
	}
	return umember
}

func GetMenberStatus(platform, game int) []MemberUpload {
	var umenber []MemberUpload
	if err := DBEngine.Db.Where("type = '5' and remark0 = ?", game).Find(&umenber).Error; err != nil {
		util.Log.Error(err)
	}
	if len(umenber) == 0 {
		if err := DBEngine.Db.Where("type between '4' and '5' and remark0=?", game).Delete(MemberUpload{}).Error; err != nil {
			util.Log.Error(err)
		}
	}
	return umenber
}
