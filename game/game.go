package game

//IGame 游戏服务接口
type IGame interface {
	EnterRoom(*Player, int, int)
	LeaveRoom(*Player, bool, int)

	DoBet(*Player, string, int, int, int, []byte) byte
	CleanOrder(int)

	FindTable(int) *Table
	FindTables() map[int]*Table

	UpdatePlayerss(int, []*PlayerStatus, bool)
	SyncHistory([]string, int)

	SendHallMsg([]byte)
	SendStatusToAll(int, int)
	SendGroupMsg([]byte, int)
	SendPlayerMsg(int)

	RemoveMoni(*Player, int)
	AddMoni(*Player, int)

	Start(int, string, string) byte
	Stop(int)

	DoResult(int, string, string) byte
	DoChangeResult(int, string, string) byte
	DoSyncResult(int, string, string) byte
}

//IClient 客户端接口
type IClient interface {
	GetGame() IGame
	InTable(int, int) byte
	LeaveTable(int) byte
}
