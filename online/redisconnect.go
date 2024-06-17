package online

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	ConnectToRedis()
}

func ConnectToRedis() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Error connecting to Redis:", err)
		return
	}

	log.Println("***********************CONNECTED TO REDIS********************* ")

	RedisClient = rdb
}
