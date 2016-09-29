package game

import (
	"go/cmkj_server_go/models"

	"strconv"
	"sync"
	"time"
)

//BettingOrder 下注单
type BettingOrder struct {
	P                  *Player
	ID                 int64
	UserID             int64
	UserName           string
	Path               string    //会员层级路径
	GameRecordID       int64     //游戏结果ID
	TableID            int       //桌号
	Stage              int       //靴数
	Inning             int       //局数
	GameNameID         int       //游戏名ID
	GameKind           int       //游戏种类ID
	GameBettingKind    int       //会员投注种类ID
	GameBettingContent string    //用户投注内容,轮盘骰子
	GameResult         string    //游戏结果
	GameResultContent  float64   //用户结果内容
	ResultType         int       //输 = 1,赢 = 2,和 = 3
	OrderNumber        string    //
	BettingAmount      float64   //投注金额
	CompensateRate     float64   //赔率
	WinLoseAmount      float64   //输赢金额
	Balance            float64   //余额
	IP                 string    //ip
	AddTime            time.Time //投注时间
	PlatformID         int       // 平台ID
	StartTime          time.Time //开始时间
}

//NewBettingOrder 新建BettingOrder
func NewBettingOrder() *BettingOrder {
	return &BettingOrder{}
}

//BettingOrders 下注单表
type BettingOrders struct {
	mutex  sync.Mutex
	orders []*BettingOrder
}

//NewBettingOrders 新建下注单表
func NewBettingOrders() *BettingOrders {
	return &BettingOrders{
		orders: make([]*BettingOrder, 0),
	}
}

//OrderConvertToTmp 下注单转换成数据库临时单结构
func OrderConvertToTmp(order *BettingOrder) *models.VideoBettingTmp {
	tmp := &models.VideoBettingTmp{}
	tmp.UserID = order.UserID
	tmp.UserName = order.UserName
	tmp.Stage = order.Stage
	tmp.Path = order.Path
	tmp.GameRecordID = strconv.FormatInt(order.GameRecordID, 10)
	tmp.TableID = strconv.Itoa(order.TableID)
	tmp.OrderNumber = order.OrderNumber
	tmp.Inning = order.Inning
	tmp.GameResult = order.GameResult
	tmp.ResultType = order.ResultType
	tmp.BettingAmount = order.BettingAmount
	tmp.CompensateRate = order.CompensateRate
	tmp.WinLoseAmount = order.WinLoseAmount
	tmp.Balance = order.Balance
	tmp.IP = order.IP
	tmp.AddTime = order.AddTime
	tmp.StartBetTime = order.StartTime
	return tmp
}

//TmpToBettingOrder 数据库临时单转成下注单
func TmpToBettingOrder(tmp models.VideoBettingTmp) *BettingOrder {
	order := &BettingOrder{}
	order.ID = tmp.ID
	order.UserID = tmp.UserID
	order.UserName = tmp.UserName
	order.Path = tmp.Path
	rid, _ := strconv.ParseInt(tmp.GameRecordID, 10, 64)
	order.GameRecordID = rid
	tid, _ := strconv.Atoi(tmp.TableID)
	order.TableID = tid
	order.Stage = tmp.Stage
	order.Inning = tmp.Inning
	order.GameNameID = tmp.GameNameID
	order.GameKind = tmp.GameKind
	order.GameBettingKind = tmp.GameBettingKind
	order.GameBettingContent = tmp.GameBettingContent
	order.GameResult = tmp.GameResult
	order.OrderNumber = tmp.OrderNumber
	order.ResultType = tmp.ResultType
	order.BettingAmount = tmp.BettingAmount
	order.CompensateRate = tmp.CompensateRate
	order.WinLoseAmount = tmp.WinLoseAmount
	order.Balance = tmp.Balance
	order.IP = tmp.IP
	order.AddTime = tmp.AddTime
	order.PlatformID = tmp.PlatformID
	order.StartTime = tmp.StartBetTime

	return order
}

//OrderToBettingOrder 数据库下单记录转成下注单
func OrderToBettingOrder(tmp models.VideoBetting) *BettingOrder {
	order := &BettingOrder{}
	order.ID = tmp.ID
	order.UserID = tmp.UserID
	order.UserName = tmp.UserName
	order.Path = tmp.Path
	rid, _ := strconv.ParseInt(tmp.GameRecordID, 10, 64)
	order.GameRecordID = rid
	tid, _ := strconv.Atoi(tmp.TableID)
	order.TableID = tid
	order.Stage = tmp.Stage
	order.Inning = tmp.Inning
	order.GameNameID = tmp.GameNameID
	order.GameKind = tmp.GameKind
	order.GameBettingKind = tmp.GameBettingKind
	order.GameBettingContent = tmp.GameBettingContent
	order.GameResult = tmp.GameResult
	order.OrderNumber = tmp.OrderNumber
	order.ResultType = tmp.ResultType
	order.BettingAmount = tmp.BettingAmount
	order.CompensateRate = tmp.CompensateRate
	order.WinLoseAmount = tmp.WinLoseAmount
	order.Balance = tmp.Balance
	order.IP = tmp.IP
	order.AddTime = tmp.AddTime
	order.PlatformID = tmp.PlatformID
	return order
}

//ToVideoBetting 下注单转成数据库下单揭露结构
func ToVideoBetting(order *BettingOrder) *models.VideoBetting {
	return &models.VideoBetting{
		ID:                 order.ID,
		UserID:             order.UserID,
		UserName:           order.UserName,
		Path:               order.Path,
		GameRecordID:       strconv.FormatInt(order.GameRecordID, 10),
		OrderNumber:        order.OrderNumber,
		TableID:            strconv.Itoa(order.TableID),
		Stage:              order.Stage,
		Inning:             order.Inning,
		GameNameID:         order.GameNameID,
		GameKind:           order.GameKind,
		GameBettingKind:    order.GameBettingKind,
		GameBettingContent: order.GameBettingContent,
		GameResult:         order.GameResult,
		ResultType:         order.ResultType,
		BettingAmount:      order.BettingAmount,
		CompensateRate:     order.CompensateRate,
		WinLoseAmount:      order.WinLoseAmount,
		Balance:            order.Balance,
		IP:                 order.IP,
		AddTime:            order.AddTime,
		PlatformID:         order.PlatformID,
	}
}

//Add 添加注单
func (orders *BettingOrders) Add(order *BettingOrder) {
	orders.mutex.Lock()
	defer orders.mutex.Unlock()
	orders.orders = append(orders.orders, order)
}

//Len 注单数量
func (orders *BettingOrders) Len() int {
	return len(orders.orders)
}

//Clear 清空注单
func (orders *BettingOrders) Clear() {
	orders.mutex.Lock()
	defer orders.mutex.Unlock()
	orders.orders = make([]*BettingOrder, 0)
}

//Orders 获取所有下注单
func (orders *BettingOrders) Orders() []*BettingOrder {
	o := make([]*BettingOrder, len(orders.orders))
	copy(o, orders.orders)
	return o
}

//LoadTmpOrders 读取临时注单
func LoadTmpOrders(gameid, where, tableid, stage, inning int) *BettingOrders {
	bettingorders := NewBettingOrders()
	if where == 1 {
		orders := models.ReadOrders(gameid, tableid, stage, inning)
		for _, o := range orders {
			bettingorders.Add(OrderToBettingOrder(o))
		}
	} else {
		orders := models.ReadTmpOrders(gameid, tableid, stage, inning)
		for _, o := range orders {
			bettingorders.Add(TmpToBettingOrder(o))
		}
	}
	return bettingorders
}

//ToOrderAll 记录注单
func ToOrderAll(orders []*BettingOrder) {
	for _, order := range orders {
		ToVideoBetting(order).Create()
	}
}

//UpdateOrderAll 更新数据库注单
func UpdateOrderAll(orders []*BettingOrder) {
	if len(orders) <= 0 {
		return
	}
	for _, bo := range orders {
		vbo := ToVideoBetting(bo)
		vbo.UpdateOrderHis()
	}
}
