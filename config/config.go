package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort            string
	MongoURI              string
	MongoDB               string
	RedisAddr             string
	RedisPass             string
	MaxConcurrentRequests int
}

func LoadConfig() *Config {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	if !strings.HasPrefix(mongoURI, "mongodb://") && !strings.HasPrefix(mongoURI, "mongodb+srv://") {
		mongoURI = "mongodb://" + mongoURI
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "booking"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPass := os.Getenv("REDIS_PASS")

	maxReqStr := os.Getenv("MAX_CONCURRENT_REQUESTS")

	maxReq, err := strconv.Atoi(maxReqStr)
	if err != nil {
		maxReq = 100
	}

	return &Config{
		ServerPort:            serverPort,
		MongoURI:              mongoURI,
		MongoDB:               mongoDB,
		RedisAddr:             redisAddr,
		RedisPass:             redisPass,
		MaxConcurrentRequests: maxReq,
	}
}
