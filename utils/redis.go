package utils

import (
	"chess_server/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisAddr,
		Password: config.Config.RedisPass,
		DB:       config.Config.RedisDB,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	if config.Config.Debug {
		if err := RDB.Del(Ctx, "players_q", "players_q_set").Err(); err != nil {
			log.Printf("Failed to clear player queues: %v", err)
		} else {
			log.Println("Redis queues 'players_q' and 'players_q_set' cleared (debug mode)")
		}
	}
}
