package models

import (
	"go/cmkj_server_go/util"
)

//BetChipLimitVideo ...
type BetChipLimitVideo struct {
	ID           int    `gorm:"column:id"`
	LowerLimit   int    `gorm:"column:LowerLimit"`
	UpperLimit   int    `gorm:"column:UpperLimit"`
	ChipGroup    string `gorm:"colimn:ChipGroup"`
	ChipOptional string `gorm:"column:ChipOptional"`
}

func (BetChipLimitVideo) TableName() string {
	return "C_BetChipLimitVideo"
}

type BetChipLimitRoulette struct {
	ID           int    `gorm:"column:id"`
	LowerLimit   string `gorm:"column:LowerLimit"`
	UpperLimit   string `gorm:"column:UpperLimit"`
	ChipGroup    string `gorm:"colimn:ChipGroup"`
	ChipOptional string `gorm:"column:ChipOptional"`
}

func (BetChipLimitRoulette) TableName() string {
	return "C_BetChipLimitRoulette"
}

func ReadVideoChip(id int) (*BetChipLimitVideo, byte) {
	var chip []BetChipLimitVideo
	if err := DBEngine.Db.Where("id = ?", id).Find(&chip).Error; err != nil {
		util.Log.Error(err)
		return nil, 1
	}
	return &chip[0], 0
}

func ReadRouletteChip(id int) (*BetChipLimitRoulette, byte) {
	var chip []BetChipLimitRoulette
	if err := DBEngine.Db.Where("id = ?", id).Find(&chip).Error; err != nil {
		util.Log.Error(err)
		return nil, 1
	}
	return &chip[0], 0
}
