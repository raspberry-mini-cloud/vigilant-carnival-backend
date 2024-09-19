package redis

import (
	"context"
	"fmt"
	"log"
	. "project/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)	

var RedisClient *redis.Client
var Context context.Context

func InitRedis(ctx context.Context) {

	Context = ctx

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

func AddTemperatureData(timestamp int64, temperature float64) {
	err := RedisClient.Do(Context, "TS.ADD", "temperature_readings", timestamp, temperature).Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckKeyExists(key string)(bool) {
	result := RedisClient.Exists(Context, key)

	return (result.Val() > 0)
}