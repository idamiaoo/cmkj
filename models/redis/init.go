package redis

import (
	"go/cmkj_server_go/conf"

	"time"

	"github.com/garyburd/redigo/redis"
)

//RedisPool  全局redis连接池
var RedisPool *redis.Pool

//OpenRedis 重redis连接池获取一个连接进行数据库访问
func OpenRedis() redis.Conn {
	return RedisPool.Get()
}

//CloseRedis 释放一个redis连接，终止数据库访问
func CloseRedis(conn redis.Conn) {
	conn.Close()
}

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
