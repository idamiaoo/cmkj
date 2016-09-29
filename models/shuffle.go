package models

import (
	"go/cmkj_server_go/util"

	"time"
)

type Shuffle struct {
	ID            int       `gorm:"column:ID;primary_key"`
	Shuffle       string    `gorm:"column:Shuffle"` //洗牌方式 A或B
	TableID       int       `gorm:"column:TableId"`
	Stage         int       `gorm:"column:Stage"`         //靴数
	BetCount      int       `gorm:"column:BetCount"`      //投注次数
	BettingAmount float64   `gorm:"column:BettingAmount"` //投注额
	VavlidAmount  float64   `gorm:"column:VavlidAmount"`  //有效投注
	WinLoseAmount float64   `gorm:"column:WinLoseAmount"` //输赢金额
	StartDate     time.Time `gorm:"column:StartDate"`
	EndDate       time.Time `gorm:"column:EndDate"`
}

func (Shuffle) TableName() string {
	return "C_Shuffle"
}

type ShuffleStatus struct {
	ID            int       `gorm:"column:ID;primary_key"`
	Shuffle       string    `gorm:"column:Shuffle"` //洗牌方式 A或B
	TableID       int       `gorm:"column:TableId"`
	Stage         int       `gorm:"column:Stage"`         //靴数
	BetCount      int       `gorm:"column:BetCount"`      //投注次数
	BettingAmount float64   `gorm:"column:BettingAmount"` //投注额
	VavlidAmount  float64   `gorm:"column:VavlidAmount"`  //有效投注
	WinLoseAmount float64   `gorm:"column:WinLoseAmount"` //输赢金额
	StartDate     time.Time `gorm:"column:StartDate"`
	EndDate       time.Time `gorm:"column:EndDate"`
}

func (ShuffleStatus) TableName() string {
	return "C_Shuffle_status"
}

func NewShuffle(shuffle string) *Shuffle {
	now := time.Now()
	return &Shuffle{
		Shuffle:   shuffle,
		StartDate: now,
		EndDate:   now,
	}
}

func (s *Shuffle) ToShuffle() {
	if err := DBEngine.Db.Create(s).Error; err != nil {
		util.Log.Error(err)
	}
}

func (s *Shuffle) LoadShuffle() {
	var (
		shuffles []Shuffle
	)
	util.Log.Info(s.TableID)
	if err := DBEngine.Db.Where("TableId = ?", s.TableID).Find(&shuffles).Error; err != nil {
		util.Log.Error(err)
		return
	}
	if len(shuffles) == 0 {
		s = NewShuffle("")
		status := s.ConvertToStatus()
		if err := DBEngine.Db.Create(status).Error; err != nil {
			util.Log.Error(err)
		}
	} else {
		status := shuffles[0]
		s.Stage = status.Stage
		s.Shuffle = status.Shuffle
		s.BetCount = status.BetCount
		s.BettingAmount = status.BettingAmount
		s.WinLoseAmount = status.WinLoseAmount
		s.VavlidAmount = status.VavlidAmount
		s.StartDate = status.StartDate
		s.EndDate = status.EndDate
	}
}

func (s *Shuffle) ConvertToStatus() *ShuffleStatus {
	return &ShuffleStatus{
		ID:            s.ID,
		TableID:       s.TableID,
		BetCount:      s.BetCount,
		BettingAmount: s.BettingAmount,
		VavlidAmount:  s.VavlidAmount,
		WinLoseAmount: s.WinLoseAmount,
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
	}
}

func (s *Shuffle) ToShuffleStatus() {
	status := s.ConvertToStatus()
	status.ToShuffleStatus()
}

func (s *ShuffleStatus) ToShuffleStatus() {
	if err := DBEngine.Db.Exec("update C_Shuffle_status set shuffle=?,stage=?,BetCount=?,BettingAmount=?,VavlidAmount=?,WinLoseAmount=?,StartDate=?,EndDate=? where tableId=?",
		s.Shuffle, s.Stage, s.BetCount, s.BettingAmount, s.VavlidAmount, s.WinLoseAmount, s.StartDate, s.EndDate, s.TableID).Error; err != nil {
		util.Log.Error(err)
	}
}
