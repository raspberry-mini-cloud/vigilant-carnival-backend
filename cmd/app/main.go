package main

import (
	. "project/internal/db"
	. "project/internal/facades"
	. "project/internal/redis"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	InitDB()
	InitRedis()

	router.POST("/temperature", PostTemperatureData)
	router.GET("/temperature", GetTemperatureDataInSpan)
	router.DELETE("/temperature/:id", DeleteTemperatureData)


	router.Run("0.0.0.0:5000")
}