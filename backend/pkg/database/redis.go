package database

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}

	opts := &redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0, // use default DB
	}

	rdb := redis.NewClient(opts)

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return rdb
}
