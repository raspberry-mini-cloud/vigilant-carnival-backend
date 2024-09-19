package main

import (
	. "project/internal/db"
	. "project/internal/facades"
	. "project/internal/redis"

	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	InitDB()
	InitRedis(ctx)

	defer Db.Close()
	defer RedisClient.Close()

	router.POST("/temperature", PostTemperatureData)
	router.GET("/temperature", GetTemperatureDataInSpan)
	router.DELETE("/temperature/:id", DeleteTemperatureData)


	router.Run("0.0.0.0:5000")
}