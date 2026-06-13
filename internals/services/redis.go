package services

import (
	"context"
	"log"
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
var ctx context.Context = context.Background()

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

func(rc *RedisConnection) Subscribe(channel string) (<-chan int, error) {
	subs := rc.Client.Subscribe(ctx, channel)
	listner := make(chan int)
	go func(){
		defer subs.Close()
	for{
		msg, err := subs.ReceiveMessage(ctx)
		if err != nil{
			log.Println(err)
			listner <- -1
			return
		}
		rawStatus := (msg.String())[0]
		status := int(rawStatus)
		listner <- status
	}
	}()
	return listner,nil
}

func(rc *RedisConnection) Publish(channel string, msg string){
	rc.Client.Publish(ctx, channel, msg)
}

