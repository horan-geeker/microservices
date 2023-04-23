package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"microservices/internal/config"
)

var ikRedisObj *ikRedis = nil

type ikRedis struct {
	client *redis.Client
}

func connect() *redis.Client {
	log.Println(fmt.Sprintf("%s:%d", config.AppConfig.RedisHost, config.AppConfig.RedisPort))
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDB,
	})

	_, err := conn.Ping().Result()
	if err != nil {
		panic(err)
	}
	return conn
}

func GetRedisConn() *redis.Client {
	if ikRedisObj != nil {
		return ikRedisObj.client
	}

	ikRedisObj = &ikRedis{
		client: connect(),
	}

	return ikRedisObj.client
}
