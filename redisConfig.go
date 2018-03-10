package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)


//redis连接池
var (
	RedisClient *redis.Pool
	HOST        string
	DB          int
)

func init() {
	HOST = "localhost:6379"
	DB  = 0
	// 建立连接池
	RedisClient = &redis.Pool{
		MaxIdle:     10,//最大空闲数
		MaxActive:   50,//最大活跃数
		IdleTimeout: 180 * time.Second,//超时时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", HOST)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", DB)
			return c, nil
		},
	}
}
