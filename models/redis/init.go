package redis

import (
	"go/cmkj_server_go/conf"

	"time"

	"github.com/garyburd/redigo/redis"
)

//RedisPool  全局redis连接池
var RedisPool *redis.Pool

//NewRedisPool 建立redis连接池
func NewRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 480 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db > 0 && db < 16 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

//Init 初始化redis连接池
func Init() {
	server := conf.Conf.DefaultString("redis_server", "127.0.0.1:6379")
	password := conf.Conf.DefaultString("redis_paswword", "")
	db := conf.Conf.DefaultInt("redis_db", 1)
	RedisPool = NewRedisPool(server, password, db)
}
