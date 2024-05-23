package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisClient() redis.Client {

	errDot := godotenv.Load("./.env")
	if errDot != nil {
		fmt.Println("Error loading .env file")
		panic(errDot)
	}

	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:%s/0?max_retries=2", os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT")))

	rdb := redis.NewClient(opt)

	return *rdb
}
