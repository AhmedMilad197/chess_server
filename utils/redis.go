package utils

import (
	"chess_server/config"
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisAddr,
		Password: config.Config.RedisPass,
		DB:       config.Config.RedisDB,
	})
}
