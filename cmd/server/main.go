package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/api"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
	"github.com/naghinezhad/BookingResourceSystem/internal/redis"
	"github.com/naghinezhad/BookingResourceSystem/internal/utils/logger"
)

func main() {
	runMigration := flag.Bool("migration", false, "run database migrations and exit")
	flag.Parse()

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

	if *runMigration {
		err = database.RunMigrations(cfg.MongoURI, cfg.MongoDB)
		if err != nil {
			log.Fatal("migration error:", err)
		}

		log.Println("Migrations applied")
		return
	}

	// redis connection
	redisClient, err := redis.NewRedis(
		ctx,
		cfg.RedisAddr,
		cfg.RedisPass,
	)
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
