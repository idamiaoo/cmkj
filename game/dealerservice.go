package game

import (
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"

	"math"
	"strings"
	"time"
)

//DealerService 荷官
type DealerService struct {
	Dealerr     *models.Dealer
	Shufflee    *models.Shuffle
	TableID     int
	HaveShuffle bool
}

//NewDealerService ...
func NewDealerService(gameid, tableid int, haveShuffle bool) *DealerService {
	tid := gameToID(gameid, tableid)
	util.Log.Debug(gameid, tid)
	s := &DealerService{
		TableID: tid,
		Dealerr: &models.Dealer{
			TableID: tid,
		},
		Shufflee: &models.Shuffle{
			TableID: tid,
		},
		HaveShuffle: haveShuffle,
	}
	s.LoadDealer()
	if haveShuffle {
		s.LoadShuffle()
	}
	return s
}

func gameToID(gameID, tableID int) int {
	switch gameID {
	case 11:
		if tableID > 3 {
			return tableID + 1
		}
		return tableID
	case 12:
		return 9 + tableID
	case 13:
		return 15
	case 14:
		return 12
	case 16:
		return 13
	default:
	}
	return 0
}

//AddWinLose 添加输赢记录
func (s *DealerService) AddWinLose(betCount int, bettingAmount, winLoseAmount, vavlid float64) {
	vavlid = math.Abs(winLoseAmount)
	if vavlid < bettingAmount {
		vavlid = bettingAmount
	}
	s.Dealerr.BetCount += betCount
	s.Dealerr.BettingAmount += bettingAmount
	s.Dealerr.WinLoseAmount -= winLoseAmount
	s.Dealerr.VavlidAmount += vavlid
	if !s.HaveShuffle {
		return
	}
	s.Shufflee.BetCount += betCount
	s.Shufflee.BettingAmount += bettingAmount
	s.Shufflee.WinLoseAmount -= winLoseAmount
	s.Shufflee.VavlidAmount += vavlid
}

// MannageIn 荷官登录
func (s *DealerService) MannageIn(name, date string) {
	util.Log.Debug(name, date)
	if strings.EqualFold(name, s.Dealerr.UserName) {
		return
	}
	s.Dealerr.EndDate = time.Now()
	s.Dealerr.ToDealer()
	var d time.Time
	if date == "" {
		d = time.Now()
	} else {
		d, _ = time.Parse(util.Layout, date)
	}
	s.Dealerr.Init(name, d)
	//s.Dealerr = models.NewDealer(name, d)
	s.ToDealerStatus()
}

//Shuffle ...
func (s *DealerService) Shuffle(t string, stage int) {
	if !s.HaveShuffle || s.Shufflee.Stage == stage {
		return
	}
	s.Shufflee.EndDate = time.Now()
	s.Shufflee.ToShuffle()
	s.Shufflee = models.NewShuffle(t)
	s.Shufflee.Stage = stage
	s.ToShuffleStatus()
}

//LoadDealer ...
func (s *DealerService) LoadDealer() {
	s.Dealerr.LoadDealer()
}

//LoadShuffle ...
func (s *DealerService) LoadShuffle() {
	s.Shufflee.LoadShuffle()
}

//ToDealerStatus ...
func (s *DealerService) ToDealerStatus() {
	//util.Log.Debug(s.Dealerr.TableID)
	s.Dealerr.ToDealerStatus()
}

//ToShuffleStatus ...
func (s *DealerService) ToShuffleStatus() {
	s.Shufflee.ToShuffleStatus()
}
