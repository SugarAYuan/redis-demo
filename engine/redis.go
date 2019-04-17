package engine

import (
	"github.com/go-redis/redis"
	"redisdemo/tools"
)

var RedisCli *redis.Client

func InitRedisEngine() {
	options := &redis.Options{
		Addr:     tools.Config.GetString("redis.host"),
		Password: tools.Config.GetString("redis.password"),
		DB:       tools.Config.GetInt("redis.db"),
		PoolSize: tools.Config.GetInt("redis.max_idle"),
	}
	RedisCli = redis.NewClient(options)
	if _, err := RedisCli.Ping().Result(); err != nil {
		panic(err)
	}

}
