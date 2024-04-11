package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	DefaultAddr     = "127.0.0.1:6379"
	DefaultPassword = ""
	DefaultDB       = 0
)

func GenerateRedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = DefaultAddr
	}
	// logger.Printf("Redis Addr: %s", redisAddr)

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = DefaultPassword
	}
	// logger.Printf("Redis Password: %s", redisPassword)

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := DefaultDB
	if redisDBStr != "" {
		redisDBTemp, err := strconv.Atoi(redisDBStr)
		if err == nil {
			redisDB = redisDBTemp
		}
	}
	// logger.Printf("Redis DB: %d", redisDB)

	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
}

func GetTopoSpecFromRedis(ctx context.Context, redisClient *redis.Client, name string) (RedisTopologySpec, error) {
	redisTopoSpec := &RedisTopologySpec{}

	specJSON, err := redisClient.Get(ctx, "cni_"+name+"_spec").Result()
	if err != redis.Nil {
		if err = json.Unmarshal([]byte(specJSON), redisTopoSpec); err == nil {
			return *redisTopoSpec, nil
		}
	}

	return *redisTopoSpec, fmt.Errorf("failed to get pod %s spec from redis", name)
}

func GetTopoStatusFromRedis(ctx context.Context, redisClient *redis.Client, name string) (RedisTopologyStatus, error) {
	redisTopoStatus := &RedisTopologyStatus{}

	statusJSON, err := redisClient.Get(ctx, "cni_"+name+"_status").Result()
	if err != redis.Nil {
		if err = json.Unmarshal([]byte(statusJSON), redisTopoStatus); err == nil {
			return *redisTopoStatus, nil
		}
	}

	return *redisTopoStatus, fmt.Errorf("failed to get pod %s status from redis", name)
}
