package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(url, password string, redisDb int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            url,
		Password:        password,
		DB:              redisDb,
		MinRetryBackoff: time.Millisecond * 100,
		DialTimeout:     time.Second * 15,
	})

	return rdb, rdb.Ping(context.TODO()).Err()
}
