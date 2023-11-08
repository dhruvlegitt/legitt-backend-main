package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisSession struct {
	Email    string
	Location string
	Exp      int64
}

func (r RedisSession) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

// Thread Safe connection pool
var RedisClient *redis.Client
var ctx = context.TODO()

func InitRedis() *redis.Client {

	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	RedisClient = redis.NewClient(options)

	pong, err := RedisClient.Ping(ctx).Result()

	if err != nil || pong != "PONG" {
		log.Fatal("")
	}

	fmt.Print("Connected to redis successfully \n")

	return RedisClient
}

func SetRedisValue(key string, value RedisSession, expireDuration time.Duration) error {
	_, err := RedisClient.Set(ctx, key, value, expireDuration).Result()
	return err
}

func GetRedisValue(key string) (string, error) {
	res, err := RedisClient.Get(ctx, key).Result()

	if err != nil {
		fmt.Printf("Key dosen't exist for %s", key)
		return "", err
	}

	return res, nil
}

func DeleteRedisValue(key string) error {
	_, err := RedisClient.Del(ctx, key).Result()
	return err
}
