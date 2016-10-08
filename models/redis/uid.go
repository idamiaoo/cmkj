package redis

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

//uid保存键
const (
	UIDKey = "cmkj_server_uid"
)

//ReloadUID 系统重启时重载uid
func ReloadUID() (int64, error) {
	rc := OpenRedis()
	defer CloseRedis(rc)
	suid, err := redis.String(rc.Do("GET", UIDKey))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	if err == redis.ErrNil {
		return 0, nil
	}
	return strconv.ParseInt(suid, 10, 64)
}

//SaveUID 保存自增uid
func SaveUID(uid string) error {
	rc := OpenRedis()
	defer CloseRedis(rc)
	_, err := rc.Do("SET", UIDKey, uid)
	return err
}
