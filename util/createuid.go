package util

import (
	"sync"
)

type Uid struct {
	uid   int64
	mutex sync.Mutex
}

var UidGenerator *Uid

func NewUid() *Uid {
	return &Uid{}
}

func (u *Uid) GetId(game int) int64 {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	n := int64(game) * 10000000000
	u.uid += 1
	if u.uid > 9999999999 {
		u.uid = 1
	}
	return u.uid + n
}

func init() {
	UidGenerator = NewUid()
}
