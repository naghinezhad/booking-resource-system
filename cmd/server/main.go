package main

import (
	"context"
	"log"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/api"
	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
	"github.com/naghinezhad/BookingResourceSystem/internal/logger"
)

func main() {
	// init logger
	logger.Init()

	// load config
	cfg := config.LoadConfig()

	// context for connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongo connection
	mongo, err := database.NewMongo(
		ctx,
		cfg.MongoURI,
		cfg.MongoDB,
	)

	if err != nil {
		log.Fatal("mongo connection error:", err)
	}

	log.Println("Mongo connected:", mongo.DB.Name())

	// run migrations
	err = database.RunMigrations(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("migration error:", err)
	}

	log.Println("Migrations applied")

	// redis connection
	redisClient := cache.NewRedis(
		cfg.RedisAddr,
		cfg.RedisPass,
	)

	err = redisClient.Ping(ctx)
	if err != nil {
		log.Fatal("redis connection error:", err)
	}

	log.Println("Redis connected")

	// setup router
	router := api.SetupRouter(
		mongo,
		redisClient,
		cfg,
	)

	log.Println("server running on port", cfg.ServerPort)

	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
