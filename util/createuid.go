package util

import (
	"go/cmkj_server_go/models/redis"

	"fmt"
	"sync"
)

//UID id生成器
type UID struct {
	uid   int64
	mutex sync.Mutex
}

//UIDGenerator 全局id生成器
var UIDGenerator *UID

//NewUID 新建一个id生成器
func NewUID() *UID {
	return &UID{}
}

//GetID 获取自增id
func (u *UID) GetID(game int) int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	n := int64(game) * 10000000000
	u.uid++
	if u.uid > 9999999999 {
		u.uid = 1
	}
	suid := fmt.Sprintf("%d", u.uid)
	redis.SaveUID(suid)
	return u.uid + n
}

//UIDInit 初始化ID生成器
func UIDInit() {
	UIDGenerator = NewUID()
	if err := UIDGenerator.reload(); err != nil {
		Log.Fatal(err)
	}
}

func (u *UID) reload() error {
	uid, err := redis.ReloadUID()
	if err != nil {
		return err
	}
	u.uid = uid
	return nil
}
