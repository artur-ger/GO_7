package ledger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitCache инициализирует клиент Redis из REDIS_ADDR / REDIS_DB / REDIS_PASSWORD.
func InitCache() error {
	dbNumber, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return fmt.Errorf("REDIS_DB должен быть числом: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("Redis не отвечает (проверьте REDIS_ADDR): %w", err)
	}

	rdb = client
	log.Println("подключились к Redis")
	return nil
}
