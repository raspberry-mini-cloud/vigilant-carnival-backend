package redis

import (
	"github.com/redis/go-redis/v9"
	"fmt"
	"context"
	"log"
)

var RedisClient *redis.Client

func InitRedis() {
    ctx := context.Background()


	RedisClient = redis.NewClient(&redis.Options{
        Addr:	  "10.40.125.129:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	pong, err := RedisClient.Ping(ctx).Result()
	fmt.Println(pong, err)
	if err != nil {
        log.Fatalln("Redis connection was refused")
    }
    fmt.Println(pong)

}