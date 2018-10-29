package main

import (
    "github.com/go-redis/redis"
    "fmt"
    "log"
)

var ikRedisObj *ikRedis = nil

type ikRedis struct {
    client *redis.Client
}

func connect() *redis.Client {
    log.Println(fmt.Sprintf("%s:%d", config.REDIS_HOST, config.REDIS_PORT))
    conn := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", config.REDIS_HOST, config.REDIS_PORT),
        Password: config.REDIS_PASSWORD,
        DB:       config.REDIS_DB,
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
