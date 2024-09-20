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

func PopulateRedisCache(temperatures []TemperatureData) {
	for _, temperature := range temperatures {
		AddTemperatureData(temperature.Timestamp.Unix(), temperature.Temperature)
	}
}

func AddTemperatureData(timestamp int64, temperature float64) {
	err := RedisClient.Do(Context, "TS.ADD", "temperature_readings", timestamp, temperature).Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func GetAllTemperatureData()([]TemperatureData) {
	result, err := RedisClient.Do(Context, "TS.RANGE", "temperature_readings", 0, time.Now().Unix()).Result()
	if err != nil {
		log.Fatalln("Could not read values from redis")
	}

	// Type assertion to ensure result is the expected type (e.g., slice of interface{})
	entries, ok := result.([]interface{})
	if !ok {
		// Handle unexpected type
		fmt.Println("Unexpected result type")
		return nil
	}

	var temperatureData []TemperatureData

	// Process each entry
	for _, entry := range entries {
		// Assuming each entry is also a slice (timestamp, value)
		data, ok := entry.([]interface{})
		if !ok || len(data) != 2 {
			fmt.Println("Unexpected entry format")
			continue
		}

		// Extract timestamp and value
		timestamp := data[0].(int64) // Convert to int64 if necessary
		temp := data[1].(float64)   // Convert to float64 or your expected data type

		temperatureData = append(temperatureData, TemperatureData{ Timestamp: time.Unix(timestamp/1000, 0), Temperature: temp})

	}
	return temperatureData

}

func CheckKeyExists(key string)(bool) {
	result := RedisClient.Exists(Context, key)

	return (result.Val() > 0)
}