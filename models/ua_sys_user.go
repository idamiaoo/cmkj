package models

import (
	"go/cmkj_server_go/util"
)

type UASysUser struct {
	UserID            int64  `gorm:"column:UserId"`
	XfKind            int    `gorm:"column:xfkind"`
	WhKind            int    `gorm:"column:whkind"`
	Email             string `gorm:"column:Email"`
	Path              string `gorm:"column:Path"`
	PromotionOtherURL string `gorm:"column:PromotionOtherURL"`
}

func (UASysUser) TableName() string {
	return "UASysUser"
}

func ReadPreOpt(userid int64) *UASysUser {
	var opt []UASysUser
	if err := DBEngine.Db.Where("userid = ?", userid).Find(&opt).Error; err != nil {
		util.Log.Error(err)
	}
	if len(opt) <= 0 {
		return nil
	}
	return &opt[0]
}
