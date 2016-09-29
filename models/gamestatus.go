package models

import (
	"go/cmkj_server_go/util"
)

type VideoGameStatus struct {
	ID         int64   `gorm:"column:ID"`
	GameNameID int     `gorm:"column:GameNameID"`
	TableID    int     `gorm:"column:Tableid"`
	Stage      int     `gorm:"column:Stage"`
	Inning     int     `gorm:"column:Inning"`
	Status     int32   `gorm:"column:Status"`
	Times      int32   `gorm:"column:Times"`
	Results    string  `gorm:"column:Results"`
	Counts     string  `gorm:"column:Counts"`
	ResultID   int64   `gorm:"column:Resultid"`
	IsOpen     bool    `gorm:"column:IsOpen"`
	Dealer     string  `gorm:"column:Dealer"`
	WinLose    float64 `gorm:"column:winLose"`
	Limit      int     `gorm:"column:limit"`
}

func (VideoGameStatus) TableName() string {
	return "C_VideoGameStatus"
}

type VideoURL struct {
	VideoURL string `gorm:"column:VideoUrl"`
}

func (VideoURL) TableName() string {
	return "C_VideoUrl"
}

func ReadVideoURL() string {
	var url []VideoURL
	DBEngine.Db.Find(&url)
	return url[0].VideoURL
}

func ReadStageHistory(gameid, tableid int) *VideoGameStatus {
	var gamestatus []VideoGameStatus
	if err := DBEngine.Db.Where("GameNameId = ? AND tableid = ?", gameid,
		tableid).Find(&gamestatus).Error; err != nil {
		util.Log.Error(err)

	}
	if len(gamestatus) <= 0 {
		return nil
	}
	return &gamestatus[0]
}

func UpdateStageHistoryWinLose(gameid, tableid int, status int32, results, counts string, winlose float64) bool {
	gamestatus := VideoGameStatus{}
	if err := DBEngine.Db.Model(&gamestatus).Where("GameNameID = ? AND Tableid = ?", gameid, tableid).Updates(map[string]interface{}{
		"Status":  status,
		"Times":   0,
		"Counts":  counts,
		"winLose": winlose,
		"Results": results,
	}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func UpdateStageHistoryStage(gameid, tableid, stage, noSub int, status int32, results, counts string) bool {
	gamestatus := VideoGameStatus{}
	if err := DBEngine.Db.Model(&gamestatus).Where("GameNameID = ? AND Tableid = ?", gameid, tableid).Updates(map[string]interface{}{
		"Results": results,
		"Status":  status,
		"Times":   0,
		"Counts":  counts,
		"Stage":   stage,
		"Inning":  noSub,
	}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func UpdateStageHistoryDealer(gameid, tableid int, isopen bool, dealer string) bool {
	gamestatus := VideoGameStatus{}
	if err := DBEngine.Db.Model(&gamestatus).Where("GameNameID = ? AND Tableid = ?", gameid, tableid).Updates(map[string]interface{}{
		"Dealer": dealer,
		"IsOpen": isopen,
	}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func UpdateStatusID(gameid, tableid, stage, noSub int, status, times int32, resultid int64) bool {
	gamestatus := VideoGameStatus{}
	if err := DBEngine.Db.Model(&gamestatus).Where("GameNameID = ? AND Tableid = ?", gameid, tableid).Updates(map[string]interface{}{
		"Times":    times,
		"Status":   status,
		"Resultid": resultid,
		"Stage":    stage,
		"Inning":   noSub,
	}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}

func UpdateStatusTime(gameid, tableid int, status, times int32) bool {
	gamestatus := VideoGameStatus{}
	if err := DBEngine.Db.Model(&gamestatus).Where("GameNameID = ? AND Tableid = ?", gameid, tableid).Updates(map[string]interface{}{
		"Status": status,
		"Times":  times,
	}).Error; err != nil {
		util.Log.Error(err)
		return false
	}
	return true
}
