package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	IPLimit           int
	IPBlockSeconds    int
	TokenLimit        int
	TokenBlockSeconds int
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ipLimit, _ := strconv.Atoi(getEnv("IP_LIMIT", "5"))
	ipBlock, _ := strconv.Atoi(getEnv("IP_BLOCK_SECONDS", "300"))
	tokenLimit, _ := strconv.Atoi(getEnv("TOKEN_LIMIT", "100"))
	tokenBlock, _ := strconv.Atoi(getEnv("TOKEN_BLOCK_SECONDS", "300"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Port:              getEnv("PORT", "8080"),
		IPLimit:           ipLimit,
		IPBlockSeconds:    ipBlock,
		TokenLimit:        tokenLimit,
		TokenBlockSeconds: tokenBlock,
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           redisDB,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
