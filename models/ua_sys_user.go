package models

import (
	"go/cmkj_server_go/util"
)

type UASysUser struct {
	UserID int64  `gorm:"column:userid"`
	XfKind int    `gorm:"column:xfkind"`
	WhKind int    `gorm:"column:whkind"`
	Email  string `gorm:"column:email"`
}

func (UASysUser) TableName() string {
	return "UASysUser"
}

func ReadPreOpt(userid int64) *UASysUser {
	var opt []UASysUser
	if err := DBEngine.Db.Where("userid = ?", userid).Find(&opt).Error; err != nil {
		util.Log.Error(err)
	}
	return &opt[0]
}
