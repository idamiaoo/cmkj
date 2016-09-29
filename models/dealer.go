package models

import (
	"go/cmkj_server_go/util"

	"time"
)

type Dealer struct {
	ID            int       `gorm:"column:ID;primary_key"`
	UserName      string    `gorm:"column:UserName"`
	TableID       int       `gorm:"column:TableId"`
	BetCount      int       `gorm:"column:BetCount"`
	BettingAmount float64   `gorm:"column:BettingAmount"`
	VavlidAmount  float64   `gorm:"column:VavlidAmount"`
	WinLoseAmount float64   `gorm:"column:WinLoseAmount"`
	StartDate     time.Time `gorm:"column:StartDate"`
	EndDate       time.Time `gorm:"column:EndDate"`
}

func (Dealer) TableName() string {
	return "C_Dealer"
}

type DealerStatus struct {
	ID            int       `gorm:"column:ID;primary_key"`
	UserName      string    `gorm:"column:UserName"`
	TableID       int       `gorm:"column:TableId"`
	BetCount      int       `gorm:"column:BetCount"`
	BettingAmount float64   `gorm:"column:BettingAmount"`
	VavlidAmount  float64   `gorm:"column:VavlidAmount"`
	WinLoseAmount float64   `gorm:"column:WinLoseAmount"`
	StartDate     time.Time `gorm:"column:StartDate"`
	EndDate       time.Time `gorm:"column:EndDate"`
}

func (DealerStatus) TableName() string {
	return "C_Dealer_status"
}

func NewDealer(name string, date time.Time) *Dealer {
	zero := time.Time{}
	if date == zero {
		date = time.Now()
	}
	return &Dealer{
		UserName:  name,
		StartDate: date,
		EndDate:   date,
	}
}

func (d *Dealer) Init(name string, date time.Time) {
	zero := time.Time{}
	d.UserName = name
	d.BetCount = 0
	d.BettingAmount = 0
	d.VavlidAmount = 0
	d.WinLoseAmount = 0
	if date == zero {
		date = time.Now()
	}
	d.StartDate = date
	d.EndDate = d.StartDate
}

func (d *Dealer) ToDealer() {
	if err := DBEngine.Db.NewRecord(*d); !err {
		util.Log.Error("insert error ")
	}
}

func (d *Dealer) LoadDealer() {
	var (
		dealers []Dealer
	)
	util.Log.Info(d.TableID)
	if err := DBEngine.Db.Where("TableId = ?", d.TableID).Find(&dealers).Error; err != nil {
		util.Log.Error(err)
		return
	}
	if len(dealers) == 0 {
		d = NewDealer("", time.Now())
		status := d.ConvertToStatus()
		if err := DBEngine.Db.Create(status).Error; err != nil {
			util.Log.Error(err)
		}
	} else {
		status := dealers[0]
		d.UserName = status.UserName
		d.BetCount = status.BetCount
		d.BettingAmount = status.BettingAmount
		d.WinLoseAmount = status.WinLoseAmount
		d.VavlidAmount = status.VavlidAmount
		d.StartDate = status.StartDate
		d.EndDate = status.EndDate
	}
}

func (d *Dealer) ConvertToStatus() *DealerStatus {
	return &DealerStatus{
		ID:            d.ID,
		UserName:      d.UserName,
		TableID:       d.TableID,
		BetCount:      d.BetCount,
		BettingAmount: d.BettingAmount,
		VavlidAmount:  d.VavlidAmount,
		WinLoseAmount: d.WinLoseAmount,
		StartDate:     d.StartDate,
		EndDate:       d.EndDate,
	}
}

func (d *Dealer) ToDealerStatus() {
	dealerstatus := d.ConvertToStatus()
	dealerstatus.ToDealerStatus()
}

func (d *DealerStatus) ToDealerStatus() {
	util.Log.Debug(d.TableID)
	if err := DBEngine.Db.Exec("update C_Dealer_status set UserName=?,BetCount=?,BettingAmount=?,VavlidAmount=?,WinLoseAmount=?,StartDate=?,EndDate=? where tableId=?",
		d.UserName, d.BetCount, d.BettingAmount, d.VavlidAmount, d.WinLoseAmount, d.StartDate, d.EndDate, d.TableID).Error; err != nil {
		util.Log.Error(err)
	}
}
