package models

import (
	"go/cmkj_server_go/util"

	"time"
)

type VideoBetting struct {
	ID                 int64     `gorm:"column:ID"`                 //[bigint]
	UserID             int64     `gorm:"column:UserID"`             //[bigint]
	UserName           string    `gorm:"column:UserName"`           //[nvarchar]
	Path               string    `gorm:"column:Path"`               //[nvarchar]  会员层级路径
	GameRecordID       string    `gorm:"column:GameRecordID"`       //"nvarchar](30) NOT NULL, 游戏结果ID
	OrderNumber        string    `gorm:"column:OrderNumber"`        //"nvarchar](20) NOT NULL,
	TableID            string    `gorm:"column:TableID"`            //"nvarchar](4) NOT NULL,桌号
	Stage              int       `gorm:"column:Stage"`              //"smallint] NOT NULL,靴数
	Inning             int       `gorm:"column:Inning"`             //"smallint] NOT NULL,局数
	GameNameID         int       `gorm:"column:GameNameID"`         //"smallint] NOT NULL,游戏名ID
	GameKind           int       `gorm:"column:GameKind"`           //"smallint] NOT NULL,游戏种类ID
	GameBettingKind    int       `gorm:"column:GameBettingKind"`    //"smallint] NOT NULL,会员投注种类ID
	GameBettingContent string    `gorm:"column:GameBettingContent"` //"nvarchar](4000) NOT NULL,用户投注内容,轮盘骰子
	GameResult         string    `gorm:"column:GameResult"`         //"nvarchar](100) NOT NULL,游戏结果
	ResultType         int       `gorm:"column:ResultType"`         //"tinyint] NOT NULL,//输 = 1,赢 = 2,和 = 3
	BettingAmount      float64   `gorm:"column:BettingAmount"`      //"decimal](12, 2) NOT NULL,投注金额
	CompensateRate     float64   `gorm:"column:CompensateRate"`     //"decimal](6, 2) NOT NULL,赔率
	WinLoseAmount      float64   `gorm:"column:WinLoseAmount"`      //"decimal](12, 2) NOT NULL,输赢金额
	Balance            float64   `gorm:"column:Balance"`            //"decimal](12, 2) NOT NULL,余额
	IP                 string    `gorm:"column:IP"`                 //"nvarchar](16) NOT NULL,
	AddTime            time.Time `gorm:"column:AddTime"`            //"datetime] NOT NULL,投注时间
	PlatformID         int       `gorm:"column:PlatformID"`         //"smallint] NOT NULL,平台ID
	VendorID           int64     `gorm:"column:vendor_id"`          //[bigint] NULL,
}

func (VideoBetting) TableName() string {
	return "C_VideoBetting"
}

type VideoBettingTmp struct {
	ID                 int64     `gorm:"column:ID"`                 //[bigint]
	UserID             int64     `gorm:"column:UserID"`             //[bigint]
	UserName           string    `gorm:"column:UserName"`           //[nvarchar]
	Path               string    `gorm:"column:Path"`               //[nvarchar]  会员层级路径
	GameRecordID       string    `gorm:"column:GameRecordID"`       //"nvarchar](30) NOT NULL, 游戏结果ID
	OrderNumber        string    `gorm:"column:OrderNumber"`        //"nvarchar](20) NOT NULL,
	TableID            string    `gorm:"column:TableID"`            //"nvarchar](4) NOT NULL,桌号
	Stage              int       `gorm:"column:Stage"`              //"smallint] NOT NULL,靴数
	Inning             int       `gorm:"column:Inning"`             //"smallint] NOT NULL,局数
	GameNameID         int       `gorm:"column:GameNameID"`         //"smallint] NOT NULL,游戏名ID
	GameKind           int       `gorm:"column:GameKind"`           //"smallint] NOT NULL,游戏种类ID
	GameBettingKind    int       `gorm:"column:GameBettingKind"`    //"smallint] NOT NULL,会员投注种类ID
	GameBettingContent string    `gorm:"column:GameBettingContent"` //"nvarchar](4000) NOT NULL,用户投注内容,轮盘骰子
	GameResult         string    `gorm:"column:GameResult"`         //"nvarchar](100) NOT NULL,游戏结果
	ResultType         int       `gorm:"column:ResultType"`         //"tinyint] NOT NULL,//输 = 1,赢 = 2,和 = 3
	BettingAmount      float64   `gorm:"column:BettingAmount"`      //"decimal](12, 2) NOT NULL,投注金额
	CompensateRate     float64   `gorm:"column:CompensateRate"`     //"decimal](6, 2) NOT NULL,赔率
	WinLoseAmount      float64   `gorm:"column:WinLoseAmount"`      //"decimal](12, 2) NOT NULL,输赢金额
	Balance            float64   `gorm:"column:Balance"`            //"decimal](12, 2) NOT NULL,余额
	IP                 string    `gorm:"column:IP"`                 //"nvarchar](16) NOT NULL,
	StartBetTime       time.Time `gorm:"column:StartBetTime"`
	AddTime            time.Time `gorm:"column:AddTime"`    //"datetime] NOT NULL,投注时间
	PlatformID         int       `gorm:"column:PlatformID"` //"smallint] NOT NULL,平台ID
	VendorID           int64     `gorm:"column:vendor_id"`  //[bigint] NULL,
	Createdate         time.Time `gorm:"column:createdate"`
}

func (VideoBettingTmp) TableName() string {
	return "C_VideoBetting_Temp"
}

func (bet *VideoBettingTmp) Create() bool {
	if err := DBEngine.Db.NewRecord(bet); !err {
		util.Log.Error("create tmp order error ")
		return false
	}
	return true

}

func (bet *VideoBetting) Create() bool {
	if err := BetEngine.Db.NewRecord(bet); !err {
		util.Log.Error("create betorder error ")
		return false
	}
	return true
}

func ReadTmpOrders(gameid, tableid, stage, inning int) []VideoBettingTmp {
	var tmporders []VideoBettingTmp
	if err := DBEngine.Db.Where(`GameNameID=? and TableID=? and Stage=? and Inning=? 
		and PlatformID <> 100`, gameid, tableid, stage, inning).Order("ID").Find(&tmporders).Error; err != nil {
		if err != nil {
			util.Log.Error(err)
			return nil
		}

	}
	return tmporders
}

func ReadOrders(gameid, tableid, stage, inning int) []VideoBetting {
	var orders []VideoBetting
	if err := BetEngine.Db.Where(`DateDiff(day,AddTime,GetDate())=0 and GameNameID=? and TableID=?
	 	and Stage=? and Inning=? and PlatformID <> 100`, gameid, tableid, stage, inning).Order("ID").Find(&orders).Error; err != nil {
		util.Log.Error(err)
		return nil
	}
	return orders
}

func DeleteTmpWithRid(gameRecorID int64) bool {
	if err := DBEngine.Db.Where("GameRecordID = ?", gameRecorID).Delete(VideoBettingTmp{}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func DeleteTmpWithOrderNumber(orderNumber string) bool {
	if err := DBEngine.Db.Where("OrderNumber = ?", orderNumber).Delete(VideoBettingTmp{}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func DeleteTmpWithName(name string, gameRecorID int64) bool {
	if err := DBEngine.Db.Where("GameRecordID = ? and userName = ?", gameRecorID, name).Delete(VideoBettingTmp{}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}
func GetTempNotPayoff(game, tableid int, exid int64) []VideoBettingTmp {
	var tmp []VideoBettingTmp
	querysql := "id in (select max(id) from C_VideoBetting_temp where GameNameID=? and tableid=? and GameRecordID <> ?  group by GameRecordID)"
	if err := DBEngine.Db.Where(querysql, game, tableid, exid).Error; err != nil {
		util.Log.Error(err)
		return nil
	}
	return tmp
}

func (bet *VideoBetting) UpdateOrderHis() bool {
	if err := BetEngine.Db.Exec("UPDATE C_VideoBetting SET GameRecordID = ?,GameBettingContent = ? ,GameResult = ?,ResultType = ?,WinLoseAmount = ?,Balance = ?,CreateDate = getDate(),vendor_id = NEXT VALUE FOR SeqGsOrder WHERE id = ?",
		bet.GameRecordID, bet.GameBettingContent, bet.GameResult, bet.ResultType, bet.WinLoseAmount, bet.Balance, bet.ID).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}
