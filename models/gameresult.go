package models

import (
	"go/cmkj_server_go/util"

	"time"
)

type VideoGameResult struct {
	ID              int64     `gorm:"column:ID"`              //[bigint] IDENTITY(1,1) NOT NULL,
	TableID         int       `gorm:"column:TableID"`         //[int] NOT NULL,
	Stage           int       `gorm:"column:Stage"`           //[int] NOT NULL,
	Inning          int       `gorm:"column:Inning"`          //[int] NOT NULL,
	GameInformation string    `gorm:"column:GameInformation"` //[nvarchar](max) NULL,
	GameResult      string    `gorm:"column:GameResult"`      //[nvarchar](100) NULL,
	GameNameID      int       `gorm:"column:GameNameID"`      //[int] NOT NULL,
	AddTime         time.Time `gorm:"column:AddTime"`         //[datetime] NOT NULL,
	Mark            string    `gorm:"column:mark"`            //[varchar](1) NOT NULL CONSTRAINT [DF_C_VideoGameResult_mark]  DEFAULT ('0'),
	VendorID        int64     `gorm:"column:vendor_id"`       //[bigint] NULL,
}

func (VideoGameResult) TableName() string {
	return "C_VideoGameResult"
}

func WriteResult(game, tableid, stage, noSub int, starttime time.Time, result, poker string, id int64, mark int) int64 {
	if id > 0 {
		res, err := DBEngine.Db.DB().Exec("update C_VideoGameResult set GameResult=?,GameInformation=?,mark=?,AddTime=? where ID=?",
			result, poker, mark, starttime, id)
		if err != nil {
			util.Log.Error(err)
			return 0
		}
		rows, _ := res.RowsAffected()
		return rows
	}
	//util.Log.Debug(starttime.Format(util.Layout))
	res, err := DBEngine.Db.DB().Exec("INSERT INTO C_VideoGameResult (TableID,Stage,Inning,GameInformation,GameResult,GameNameID,AddTime, mark) VALUES(?,?,?,?,?,?,?,?)",
		tableid, stage, noSub, poker, result, game, starttime, mark)
	if err != nil {
		util.Log.Error(err)
		return 0
	}
	rows, _ := res.LastInsertId()
	return rows
}

func ReadResult(game, tableid, stage int) ([]VideoGameResult, error) {
	var results []VideoGameResult
	querysql := "select * From [C_VideoGameResult] t1," +
		" ( select top 1 stage,AddTime From [C_VideoGameResult] t2 where GameNameID=? and TableID=? order by id desc) t2 " +
		"where t2.stage=? and GameNameID=? and TableID=? and t1.Stage = t2.Stage and  DATEDIFF(n,t1.addtime,t2.AddTime)<240  order by id"

	rows, err := DBEngine.Db.Raw(querysql, game, tableid, stage, game, tableid).Rows()
	if err != nil {
		util.Log.Error(err)
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var re VideoGameResult
		if err := DBEngine.Db.ScanRows(rows, &re); err != nil {
			util.Log.Error(err)
			continue
		}
		results = append(results, re)
	}
	return results, nil

}
