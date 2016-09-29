package main

import (
	"go/cmkj_server_go/game"
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
)

//Pay 下注金额变动详情
type Pay struct {
	WinLose float64
	Total   float64
	PayStr  string
	Pay     float64
	Rebate  float64
}

//NewPay ...
func NewPay() *Pay {
	return &Pay{}
}

//Init ...
func (p *Pay) Init() {
	p.WinLose = 0
	p.Total = 0
	p.PayStr = ""
	p.Pay = 0
	p.Rebate = 0
}

func converBetToSlice(betstr string) []float64 {
	betstrTmp := strings.TrimSuffix(betstr, "^")
	bets := strings.Split(betstrTmp, "^")
	bet := make([]float64, len(bets))
	for i, s := range bets {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			util.Log.Error(err)
		}
		bet[i] = f
	}
	return bet
}

func converBetToStr(bet []float64) string {
	bs := ""
	if bet == nil {
		return bs
	}
	for _, f := range bet {
		if f > 0 {
			bs += strconv.FormatFloat(f, 'f', -1, 64)
		}
		bs += "^"
	}
	return bs
}

func joinBetWithSlice(bet []float64, bet2 float64, index int) []float64 {
	if bet == nil {
		bet = make([]float64, GBetmaxid)
	}
	bet[index-1] += bet2
	return bet
}

func joinBetwithStr(bet1 string, bet2 float64, index int) string {
	bet := make([]float64, GBetmaxid)
	betTmp := strings.TrimSuffix(bet1, "^")
	bets := strings.Split(betTmp, "^")
	for i, s := range bets {
		f, _ := strconv.ParseFloat(s, 64)
		bet[i] = f
	}
	bs := ""
	bet[index-1] += bet2
	for _, f := range bet {
		if f > 0 {
			bs += strconv.FormatFloat(f, 'f', -1, 64)
		}
		bs += "^"
	}
	return bs
}

func joinBet(bet1, bet2 string) string {
	if bet1 == "" {
		return bet2
	}
	bet1Tmp := strings.TrimSuffix(bet1, "^")
	bet2Tmp := strings.TrimSuffix(bet2, "^")
	bets1 := strings.Split(bet1Tmp, "^")
	bets2 := strings.Split(bet2Tmp, "^")

	if len(bets1) != len(bets2) {
		return bet2
	}
	sb := ""
	bet := make([]float64, len(bets2))
	for i, s := range bets2 {
		b1, _ := strconv.ParseFloat(s, 64)
		b2, _ := strconv.ParseFloat(bets1[i], 64)
		bet[i] = b1 + b2
	}
	for _, f := range bet {
		sb += strconv.FormatFloat(f, 'f', -1, 64)
		sb += "^"
	}
	return sb
}

func (p *Pay) getPayBool(result []byte, bet float64, index, t int, rate []float64) bool {
	var pay float64
	if index >= len(result) || index < 0 {
		return false
	}
	util.Log.Debug(result, bet, index, t, rate)
	if result[index] > 0 {
		if index == 23 || index == 24 {
			pay = bet * float64(result[index])
		} else {
			pay = bet*rate[index] + bet
		}

		if index == game.Z {
			if t == GMianyong || t == GLianhuanmianyong {
				if result[25] == 1 {
					pay = bet*0.5 + bet
				} else {
					pay = bet*1 + bet
				}
			}
		}
		p.Pay += pay

	} else if result[game.X] == 1 {
		if index == game.Z || index == game.X {
			pay = bet
			p.Rebate -= bet
			p.Pay += pay
		}
	}

	p.Total += bet
	p.Rebate += bet
	p.WinLose = p.Pay - p.Total
	p.Pay = util.Round(p.Pay, 2)
	p.WinLose = util.Round(p.WinLose, 2)
	return true
}

func (p *Pay) getPay(result []byte, betstr string, rate []float64, rebateMode int) {
	p.Init()
	bet := converBetToSlice(betstr)

	if len(bet) > len(rate) {
		return
	}
	bets := make([]float64, len(bet))
	for i := 0; i < len(bet); i++ {
		p.getPayBool(result, bet[i], i, rebateMode, rate)
	}

	bs := ""
	p.WinLose = p.Pay - p.Total
	for _, f := range bets {
		bs += strconv.FormatFloat(f, 'f', -1, 64)
		bs += "^"
	}
	p.PayStr = bs
}
