package utils

import (
	"admin/conf"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedisClient() {
	redisAddr := fmt.Sprintf("%s:%d", conf.Config.Redis.Address, conf.Config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: conf.Config.Redis.Password,
		DB:       0,
	})
	RedisClient = rdb
}

func InitTestRedisClient(s *miniredis.Miniredis) {
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	RedisClient = rdb
}
