package game

import (
	"go/cmkj_server_go/models"
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
)

//Bet 下注信息
type Bet struct {
	Betstr   string
	BetData  []float64
	TotalBet float64
	Pay      float64
	WinLost  float64
	WinLost2 float64
	IsBet    bool
	Last     int64
}

//NewBet 获取下注实例
func NewBet() *Bet {
	return &Bet{}
}

//Player 玩家信息
type Player struct {
	ID                int64
	Name              string
	Balance           float64
	Bets              map[int]*Bet
	Bettime           int64
	Online            bool
	Game              int
	Home              int
	Vt                int
	Seat              int
	Limit             int
	LimitMin          int
	Maxprofit         float64
	ClassID           int
	IsLock            int
	Xmmod             string
	Moneysort         string
	PreSequence       string
	LoginTime         string
	IP                string
	Platform          int
	Type              int
	IsTips            int
	Tips              float64
	Xy                string
	Session           string
	Code              int
	Send              chan []byte
	PromotionOtherURL string
}

//NewPlayer 新建玩家信息实例
func NewPlayer() *Player {
	return &Player{}
}

//GetLimit 获取限压详情
func (p *Player) GetLimit(ids string) byte {
	ids = strings.TrimRight(ids, ",")
	strs := strings.Split(ids, ",")
	id, _ := strconv.Atoi(strs[0])
	if p.Game == util.Roulette {
		chip, re := models.ReadRouletteChip(id)
		if re != 0 {
			return re
		}
		min, _ := strconv.ParseFloat(chip.LowerLimit, 64)
		max, _ := strconv.ParseFloat(chip.UpperLimit, 64)
		p.LimitMin = int(min)
		p.Limit = int(max)
	} else {
		chip, re := models.ReadVideoChip(id)
		if re != 0 {
			return re
		}
		p.LimitMin = chip.LowerLimit
		p.Limit = chip.UpperLimit
	}
	return 0
}

//SendMsgToPlayer  发送消息给玩家个人
func (p *Player) SendMsgToPlayer(msg []byte) {
	p.Send <- msg
}
