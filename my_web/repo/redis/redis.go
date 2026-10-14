package redis

// require localhost.com/xuexi v0.0.0
// replace localhost.com/xuexi => E:/work/goproject/GoPath/src/xuexi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	rediscli "github.com/go-redis/redis/v8"
)

// redis连接池
var pool *redis.Pool

func init() {
	// 实例化一个连接池
	pool = &redis.Pool{
		MaxIdle:     16,  // 最初的连接数量
		MaxActive:   0,   // 最大的连接数量，不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, // 连接关闭时间 300 秒（300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { // 要连接的redis数据库
			return redis.Dial("tcp", "localhost:6379")
		},
		// Dial: func() (redis.Conn, error) {
		//     c, err := redis.Dial("tcp", "localhost:6379")
		//     if err != nil {
		//         return nil, err
		//     }
		//     if _, err := c.Do("AUTH", "123456"); err != nil {
		//         c.Close()
		//         return nil, err
		//     }
		//     return c, nil
		// },
	}
}

type RedisInter interface {
	SetHashInterface(x interface{}, key string) error
	GetHashValue(k, field string) (string, error)
	GetHashAllValue(hashKey string, bDelete bool) map[string]string
	SetOneHashData(key, field string, x interface{}) error
	GetValue(k string) (string, error)
	SetValueByTTL(k string, v any, ttl time.Duration) error
	SetValue(k string, v any) error
	SetOneHashDate(key, field string, x interface{}) error
}

type redisReal struct{}

func GetRedisInter() *redisReal {
	return &redisReal{}
}

func (rd *redisReal) SetHashInterface(x interface{}, key string) error {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去
	m, ok := x.(map[string]interface{})
	if ok {
		for k, v := range m {
			_, err := c.Do("HSet", key, k, v)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func (rd *redisReal) SetOneHashData(key, field string, x interface{}) error {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去

	_, err := c.Do("HSet", key, field, x)
	if err != nil {
		for i := 0; i < 3; i++ { // 失败重试
			_, err = c.Do("HSet", key, field, x)
			if err != nil {
				fmt.Printf("HSet error：%s, %s, %s\n", key, field, err.Error())
			} else {
				break
			}
		}
		return err
	}
	return nil
}

func (rd *redisReal) GetHashValue(k, field string) (string, error) {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去
	reply, err := c.Do("HGet", k, field)
	if err != nil || reply == nil {
		return "", err
	}
	result := ""
	switch v := reply.(type) {
	case uint8:
		result = strconv.Itoa(int(v))
	case []uint8:
		result = string(v)
	case int:
		result = strconv.Itoa(int(v))
	case string:
		result = v
	}
	return result, err
}

func (rd *redisReal) GetHashAllValue(hashKey string, bDelete bool) map[string]string {
	// 创建 Redis 客户端
	rdb := rediscli.NewClient(&rediscli.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // 密码（如果没有密码则为空）
		DB:       0,                // 默认数据库
	})

	// 创建上下文
	ctx := context.Background()

	// 使用 HGetAll 获取哈希表的所有字段和值
	result, err := rdb.HGetAll(ctx, hashKey).Result()
	if err != nil {
		fmt.Printf("Failed to get all fields from hash: %v\n", err)
	}

	// 要删除的字段
	fieldsToDelete := []string{}
	for k, v := range result {
		if v == "" {
			fieldsToDelete = append(fieldsToDelete, k)
		}
	}
	if bDelete && len(fieldsToDelete) > 0 {
		// 使用 HDel 删除指定的字段
		_, err = rdb.HDel(ctx, hashKey, fieldsToDelete...).Result()
		if err != nil {
			fmt.Printf("Failed to delete fields from hash: %v\n", err)
		}

		// 使用 HDel 删除指定的字段
		_, err = rdb.HDel(ctx, "stklimitup", fieldsToDelete...).Result()
		if err != nil {
			fmt.Printf("Failed to delete fields from hash: %v\n", err)
		}
	}

	// 使用 HGetAll 获取哈希表的所有字段和值
	resultNew, err := rdb.HGetAll(ctx, hashKey).Result()
	if err != nil {
		fmt.Printf("Failed to get all fields from hash: %v\n", err)
	}

	return resultNew
}

func (rd *redisReal) GetValue(k string) (string, error) {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去
	reply, err := c.Do("Get", k)
	if err != nil || reply == nil {
		return "", err
	}
	result := ""
	switch v := reply.(type) {
	case uint8:
		result = strconv.Itoa(int(v))
	case []uint8:
		result = string(v)
	case int:
		result = strconv.Itoa(int(v))
	case string:
		result = v
	}
	return result, err
}
func (rd *redisReal) SetValue(k string, v any) error {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去
	_, err := c.Do("set", k, v)

	return err
}

func (rd *redisReal) SetValueByTTL(k string, v any, ttl time.Duration) error {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去
	_, err := c.Do("set", k, v)
	if err != nil {
		return nil
	}
	_, err = c.Do("expire", k, ttl)

	return err
}

func (rd *redisReal) SetOneHashDate(key, field string, x interface{}) error {
	c := pool.Get() // 从连接池中取一个连接
	defer c.Close() // 函数运行结束，把连接池放回去

	_, err := c.Do("HSet", key, field, x)
	if err != nil {
		for i := 0; i < 3; i++ { // 失败重试
			_, err = c.Do("HSet", key, field, x)
			if err != nil {
				fmt.Printf("HSet error：%s, %s, %s\n", key, field, err.Error())
			} else {
				break
			}
		}
		return err
	}
	return nil
}
