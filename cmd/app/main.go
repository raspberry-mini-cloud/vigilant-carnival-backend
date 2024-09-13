package main

import (
	. "project/internal/db"
	. "project/internal/facades"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	InitDB()

	router.POST("/temperature", PostTemperatureData)
	router.GET("/temperature", GetTemperatureDataInSpan)
	router.DELETE("/temperature/:id", DeleteTemperatureData)


	router.Run("0.0.0.0:5000")
}