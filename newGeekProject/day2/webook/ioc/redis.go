package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type redisConfig struct {
		Addr string `json:"addr"`
	}
	var cfg redisConfig
	err := viper.UnmarshalKey("redis.addr", &cfg)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return redisClient
}
