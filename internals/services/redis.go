package services

import (
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct{
	Client *redis.Client
	db int
}

var redisCon *RedisConnection
var mutex sync.Mutex = sync.Mutex{}

// single connection for entire project
func FetchRedisConnection() *RedisConnection {

	mutex.Lock()
	defer mutex.Unlock()
	
	if(redisCon != nil){
		return redisCon
	}
	db := os.Getenv("REDIS_DB")
	if db == ""{
		db = "0"
	}
	rdb, _ := strconv.Atoi(db)
	host := os.Getenv("REDIS_HOST")
	if host == ""{
		host = "localhost:6379"
	}
	rc := redis.NewClient(&redis.Options{DB: rdb, Addr: host})
	redisCon = &RedisConnection{Client: rc, db: rdb}

	return redisCon
}

