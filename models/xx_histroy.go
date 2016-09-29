package models

import (
	"go/cmkj_server_go/util"

	"time"
)

type XXHistory struct {
	ID        int64     `gorm:"column:id;primary_key"`
	UserName  string    `gorm:"column:username"`
	MoneySort string    `gorm:"column:moneysort"`
	BallTime  time.Time `gorm:"column:balltime"`
	Remark0   string    `gorm:"column:remark0"`
	Remark1   string    `gorm:"column:remark1"`
	Remark2   string    `gorm:"column:remark2"`
	Tzmoney   float64   `gorm:"column:tzmoney"`
	Comm      float64   `gorm:"column:comm"`
	Total     float64   `gorm:"column:total"`
	Path      string    `gorm:"column:path"`
}

func (XXHistory) TableName() string {
	return "xx_history"
}

//Create 流水记录
func (h *XXHistory) Create() bool {
	if err := HisEngine.Db.Create(h).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}
