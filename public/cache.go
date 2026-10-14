package public

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	pwd := GetConfig("redis.pwd")
	addr := fmt.Sprintf("%s:%s", GetConfig("redis.host"), GetConfig("redis.port"))
	client := redis.NewClient(&redis.Options{
		Addr:     addr, // url
		Password: pwd,
		DB:       0, // 0号数据库
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		LOG.Errorf("redis connect ping failed, err:", err)
	} else {
		LOG.Info("redis connect ping response:", pong)
		RDB = client
	}
	return RDB
}

func GetRedis() *redis.Client {
	if RDB == nil {
		return InitRedis()
	}
	return RDB
}
