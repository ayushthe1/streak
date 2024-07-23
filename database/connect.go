package database

import (
	"log"
	"os"

	"github.com/ayushthe1/streak/models"
	"github.com/ayushthe1/streak/msgqueue"
	"github.com/ayushthe1/streak/wv"
	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	dsn := os.Getenv("DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("***********************FAIL******************* ")
		panic("Could not connect to the database")
	} else {
		log.Println("connect to postgres successfully")
	}

	log.Println("***********************CONNECTED TO POSTGRES DATABASE******************* ")
	DB = database
	database.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.ContactList{},
		&models.ActivityEvent{},
	)

	msgqueue.ConnectoRabbitMQ()

	wv.ConnectToWeaviate()

	// redisAddress := os.Getenv("REDIS_ADDRESS")
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     redisAddress,
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// _, err = rdb.Ping(ctx).Result()
	// if err != nil {
	// 	fmt.Println("Error connecting to Redis:", err)
	// 	return
	// }
	// log.Println("***********************CONNECTED TO REDIS********************* ")

	// RedisClient = rdb

}
