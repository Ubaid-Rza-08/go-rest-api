package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	goredis "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *goredis.Client {

	db := 0

	if v := os.Getenv("REDIS_DB"); v != "" {

		parsed, err := strconv.Atoi(v)

		if err == nil {
			db = parsed
		}
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	_, err := client.Ping(Ctx).Result()

	if err != nil {
		log.Fatal("REDIS CONNECTION FAILED:", err)
	}

	log.Println("REDIS CONNECTED")

	return client
}