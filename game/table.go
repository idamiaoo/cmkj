package game

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/models"
	//"go/cmkj_server_go/network"
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
	"sync"
	"time"
)

//Status 牌桌状态
type Status struct {
	Num          int
	Shoe         int       //场
	Game         int       //次
	Status       int32     //状态
	Time         int32     //倒计时
	TotalTime    int       //设置时长 //
	StartTime    time.Time //开始时间
	GameID       int64     //局编号
	Result       string    //结果 //空
	Poker        string    //扑克 //空
	PayTime      int32     //结算时延
	IsOpen       int       //
	Bet          []int     //
	VirtualBet   []int     //
	WinLose      float64   //输赢
	OlineNo      int       //在线数
	PreID        int64     //
	StatusMutext sync.Mutex
}

//NewStatus ...
func NewStatus(betAreas int) *Status {
	return &Status{
		Bet:        make([]int, betAreas),
		VirtualBet: make([]int, betAreas),
	}
}

//MakeMsg ...
func (s *Status) MakeMsg() []byte {
	//util.Log.Debug(s)
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_STATUS)
	buf.WriteByte(byte(s.Num))
	buf.WriteInt(int32(s.Shoe))
	buf.WriteInt(int32(s.Game))
	buf.WriteByte(byte(s.Status))
	buf.WriteShort(int16(s.Time))
	buf.WriteUTF(s.Poker)
	buf.WriteUTF(s.Result)
	var sb string
	for i := range s.Bet {
		sb += (strconv.Itoa(s.Bet[i] + s.VirtualBet[i])) + "^"
	}
	buf.WriteUTF(sb)
	//util.Log.Debug(sb)
	sb = ""
	for _, v := range s.Bet {
		sb += (strconv.Itoa(v)) + "^"
	}
	//util.Log.Debug(sb)
	buf.WriteUTF(sb)
	buf.WriteDouble(s.WinLose)
	return buf.Bytes()
}

//AddBet 添加下注信息
func (s *Status) AddBet(index, bet int) {
	s.StatusMutext.Lock()
	defer s.StatusMutext.Unlock()
	s.Bet[index] += bet
}

//ClearBet 清空下注信息
func (s *Status) ClearBet() {
	s.StatusMutext.Lock()
	defer s.StatusMutext.Unlock()
	for i := range s.Bet {
		s.Bet[i] = 0
	}
	for i := range s.VirtualBet {
		s.VirtualBet[i] = 0
	}
}

//Config 牌桌设置
type Config struct {
	Num         int       //
	Name        string    // 空
	Time        int32     //倒计时
	Rate        []float64 //
	VirtualBet  int
	VirtualStr  string //
	Password    string //
	Manager     *Player
	Dealer      string //
	Limit       string //
	LimitChange bool   //
}

//NewConfig ...
func NewConfig(rate int) *Config {
	return &Config{
		Rate: make([]float64, rate),
	}
}

//SetRate 设置赔率
func (c *Config) SetRate(rate string) bool {
	rate = strings.TrimSuffix(rate, "#")
	rates := strings.Split(rate, "#")
	if len(rates) < len(c.Rate) {
		util.Log.Error("赔率设置错误")
		return false
	}
	for i, v := range rates {
		r, err := strconv.ParseFloat(v, 64)
		if err != nil {
			util.Log.Error("赔率设置错误")
			return false
		}
		c.Rate[i] = r
	}
	return true
}

// LoadConfig 加载设置
func (c *Config) LoadConfig(name string) bool {
	c.Password = conf.Conf.String(conf.GetSectionkey(name, "pwd"))
	c.Time = int32(conf.Conf.DefaultInt(conf.GetSectionkey(name, "time"), 30))
	c.VirtualBet = conf.Conf.DefaultInt(conf.GetSectionkey(name, "virtualBet"), 0)
	c.VirtualStr = conf.Conf.String(conf.GetSectionkey(name, "virtualStr"))
	c.Name = conf.Conf.DefaultString(conf.GetSectionkey(name, "name"), strconv.Itoa(c.Num))
	rate := conf.Conf.String(conf.GetSectionkey(name, "rate"))
	return c.SetRate(rate)
}

// MakeMsg ...
func (c *Config) MakeMsg() []byte {
	//util.Log.Debug(*c)
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_CONFIG)
	buf.WriteByte(byte(c.Num))
	buf.WriteInt(int32(c.Time))
	buf.WriteUTF(c.Name)
	buf.WriteUTF(c.Dealer)
	buf.WriteUTF(time.Now().Format("2006-01-02 15:04:05"))
	buf.WriteUTF(c.Limit)
	return buf.Bytes()
}

//History 牌桌历史记录
type History struct {
	Num    int
	Ways   string
	Last   string //空
	Counts string //
	Poker  string //空
	Pokers string
	Tj     []int
}

//NewHistory ...
func NewHistory(tjs int) *History {
	return &History{
		Tj: make([]int, tjs),
	}
}

//MakeMsg ...
func (h *History) MakeMsg(way int) []byte {
	//util.Log.Debug(*h)
	buf := util.NewByteBuffer()
	buf.WriteShort(int16(util.P_WAY))
	buf.WriteByte(byte(h.Num))
	if way == 1 {
		buf.WriteUTF(h.Last)
		buf.WriteUTF(h.Counts)
		buf.WriteUTF(h.Poker)
	} else {
		buf.WriteUTF(h.Ways)
		buf.WriteUTF(h.Counts)
		buf.WriteUTF(h.Pokers)
	}
	return buf.Bytes()
}

//Init ...
func (h *History) Init() {
	h.Ways = ""
	h.Pokers = ""
	h.Last = "q"
	h.Counts = ""
	for i := range h.Tj {
		h.Tj[i] = 0
		h.Counts += "0^"
	}
}

func (h *History) converTongji() {
	str := strings.TrimSuffix(h.Counts, "^")
	strs := strings.Split(str, "^")
	if len(strs) < len(h.Tj) {
		h.Tj = make([]int, len(h.Tj))
	}
	for i, v := range strs {
		t, _ := strconv.Atoi(v)
		h.Tj[i] = t
	}
}

//Table 牌桌
type Table struct {
	Game        int
	Num         int
	Statu       *Status
	Conf        *Config
	His         *History
	IsWriteDb   bool
	Orders      *BettingOrders
	DealerServe *DealerService
	PlayerList  map[string]*Player
	MoniLst     map[string]*Player
	PlayerMutex sync.Mutex
	MoniMutex   sync.Mutex
	OrderMutex  sync.Mutex
}

//MakeMsg ...
func (table *Table) MakeMsg(sendway bool) []byte {
	buf := util.NewByteBuffer()
	buf.WriteShort(util.P_STATUS1)
	buf.WriteByte(byte(table.Game))
	buf.WriteByte(byte(table.Num))
	buf.WriteInt(int32(table.Statu.Shoe))
	buf.WriteInt(int32(table.Statu.Game))
	buf.WriteByte(byte(table.Statu.Status))
	util.Log.Debug(table.Statu.Time)
	buf.WriteShort(int16(table.Statu.Time))
	if sendway {
		buf.WriteUTF(table.His.Last)
		buf.WriteUTF(table.His.Counts)
	} else {
		buf.WriteUTF("")
		buf.WriteUTF("")
	}
	buf.WriteByte(byte(table.Statu.IsOpen))
	buf.WriteUTF(table.Conf.Dealer)
	if table.Conf.LimitChange {
		table.Conf.LimitChange = false
		buf.WriteUTF(table.Conf.Limit)
	} else {
		buf.WriteUTF("")
	}
	return buf.Bytes()
}

//GetStageHistory 获取游戏历史信息
func (table *Table) GetStageHistory() bool {
	his := models.ReadStageHistory(table.Game, table.Num)
	if his == nil {
		return false
	}
	table.His.Ways = his.Results
	table.His.Counts = his.Counts
	table.His.converTongji()
	table.Statu.Shoe = his.Stage
	table.Statu.Game = his.Inning
	table.Statu.Status = his.Status
	table.Statu.GameID = his.ResultID
	table.Statu.WinLose = his.WinLose
	table.Conf.Limit = strconv.Itoa(his.Limit)
	if table.Statu.Status > 0 && table.Statu.Status <= util.S_STOP {
		table.Statu.Status = util.S_STOP
	} else {
		table.Statu.Status = util.S_OVER
	}
	return true
}

//ReadGamePoker 读取扑克
func (table *Table) ReadGamePoker() bool {
	if table.Game != 11 && table.Game != 12 {
		return true
	}
	results, err := models.ReadResult(table.Game, table.Num, table.Statu.Shoe)
	if err != nil {
		return false
	}
	util.Log.Info(len(results))
	table.His.Pokers = ""
	for _, r := range results {
		table.His.Pokers += r.GameInformation + ","
	}
	util.Log.Info(table.His.Pokers)
	return true
}
