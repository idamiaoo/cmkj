package util

const (
	SYS_PORT = "PORT"
	//客户端命令
	M_LOGIN1 = 10100 //登录
	M_LOGIN  = 10101 //登录,进桌
	M_LOGIN2 = 10109 //登录

	M_BET    = 10102 //下注
	M_INROOM = 10103 //进入房间
	M_EXIT   = 10104 //退出房间
	//M_TIP     = 10105 //小费
	M_CHGNICK = 10106 //更改昵称

	M_MG    = 10901 //管理命令
	M_HEAR  = 19999 //心跳包
	M_HEAR1 = 29999 //心跳包

	//推送消息
	P_STATUS1 = 11001 //状态,大厅
	P_STATUS  = 10001 //状态,牌桌
	P_PLAYER  = 11002 //玩家，大厅
	P_PLAYER1 = 10002 //玩家，牌桌

	P_WAY       = 10003 //路单
	P_CONFIG    = 10004 //设置
	P_VIPROOM   = 10005 //包间
	P_GETOUT    = 11009 //踢出
	P_ONLINENUM = 11003 //在线人数
	P_VIDEOURL  = 11004 //video_url

	SYS_FLASHRE  = "<cross-domain-policy><allow-access-from domain='*' to-ports='*' /></cross-domain-policy>"
	SYS_FLASHGET = "<policy-file-request/>"
)

const (
	DO_ONLINE = 10 //在线
	DO_EXIT   = 0  //离线
)

const (
	//SYS_PORT  = "PORT"
	Roulette  = 13
	S_SHUFFLE = 0 //洗牌
	S_START   = 1 //开始
	S_STOP    = 2 //停止
	S_PAYOFF  = 3 //结算
	S_OVER    = 4 //结束
	S_INVALID = 5 //无效
	//M_LOGIN   = 10101 //登录
	//M_BET = 10102 //下注
	//M_INROOM  = 10103 //进入房间
	//M_EXIT    = 10104 //退出房间
	M_TIPS = 10105 //小费
	//M_MG      = 10901 //管理命令
	M_MGLOGIN = 10900 //管理登陆,1台面,2,
	//M_HEAR    = 19999 //心跳包
	//P_STATUS1 = 11001 //状态
	//P_STATUS  = 10001 //状态
	//P_PLAYER  = 10002 //玩家
	//P_WAY     = 10003 //路单
	//P_CONFIG  = 10004 //设置
	//P_VIPROOM = 10005 //包间
)
