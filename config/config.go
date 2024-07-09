package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers  string
	ConsumerGroup string
	Topic         string
	DLQTopic      string
	MaxRetries    int
	RetryDelay    time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		ConsumerGroup: os.Getenv("CONSUMER_GROUP"),
		Topic:         os.Getenv("TOPIC"),
		DLQTopic:      os.Getenv("DLQ_TOPIC"),
		MaxRetries:    getIntFromEnv("MAX_RETRIES", 5),
		RetryDelay:    getDurationFromEnv("RETRY_DELAY", 5*time.Second),
	}
}

func getIntFromEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationFromEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
