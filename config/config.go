package config

import (
	"os"
	"strconv"
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

	maxReqStr := os.Getenv("MAX_CONCURRENT_REQUESTS")

	maxReq, err := strconv.Atoi(maxReqStr)
	if err != nil {
		maxReq = 1000
	}

	return &Config{
		ServerPort:            os.Getenv("SERVER_PORT"),
		MongoURI:              os.Getenv("MONGO_URI"),
		MongoDB:               os.Getenv("MONGO_DB"),
		RedisAddr:             os.Getenv("REDIS_ADDR"),
		RedisPass:             os.Getenv("REDIS_PASS"),
		MaxConcurrentRequests: maxReq,
	}
}
