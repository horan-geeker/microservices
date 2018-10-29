package main

import (
    "github.com/jinzhu/configor"
    "github.com/joho/godotenv"
    "log"
)

var config = struct {
    APP_HOST    string `default:"127.0.0.1"`
    APP_PORT    string `default:"80"`
    APP_ENV     string `default:"local"`
    APP_VERSION string `default:"v0-1-0"`
    APP_DEBUG   string `default:"false"`

    REDIS_HOST     string `default:"127.0.0.1"`
    REDIS_PASSWORD string `default:""`
    REDIS_PORT     int    `default:"6379"`
    REDIS_DB       int    `default:"0"`
}{}

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    configor.Load(&config)
}
