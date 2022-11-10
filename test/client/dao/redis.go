package dao

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPool *redis.Pool //全局的redis连接池 用户标签数据
)

func Init() {
	redisPool = initRedisPool()
	if redisPool == nil {
		panic("redisPoolInstanceUser Err!")
	}
}

func initRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisServerAddr)
			if err != nil {
				return nil, err
			}
			if redisServerIsAuth {
				if _, err := conn.Do("AUTH", redisServerPassword); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
