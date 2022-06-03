package rediscli

import (
	"context"
	"github.com/go-redis/redis/v8"
	"orchestrator/config"
	"sync"
	"time"
)

var once sync.Once
var instance *redis.Client

func GetRedisCliInstance() *redis.Client {

	once.Do(func() {
		rCfg := config.GetConfigInstance().Client.Redis

		c := redis.NewClient(&redis.Options{
			Addr:         rCfg.Host + ":" + rCfg.Port,
			DB:           rCfg.DB,
			PoolSize:     rCfg.PoolSize,
			MinIdleConns: rCfg.MinIdleConn,
			IdleTimeout:  time.Duration(rCfg.IdleTimeout) * time.Second,
		})
		if _, err := c.Ping(context.Background()).Result(); err != nil {
			panic("redis connection refused ")
		}

		instance = c
	})

	return instance
}
