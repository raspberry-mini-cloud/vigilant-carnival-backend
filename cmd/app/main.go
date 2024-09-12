package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


var db *sql.DB

func InitDB() (*sql.DB, error) {
	postgres_password := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://postgres:%s@10.40.125.129:5432/coolkeeper_data?sslmode=disable", postgres_password)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Verify the connection is valid
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}


func main() {
	router := gin.Default()

	var err error 

	db, err = InitDB()
	if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

	router.POST("/temperature", postTemperatureData)
	router.GET("/temperature", getTemperatureDataInSpan)
	router.DELETE("/temperature/:id", deleteTemperatureData)


	router.Run("0.0.0.0:5000")
}

type TemperatureData struct {
	Temperature float64 `json:"temperature"`
	Timestamp time.Time  `json:"timestamp"`
	//Location string `json:"location"`
}

func postTemperatureData(c *gin.Context) {
	var tempData TemperatureData

	// Bind the JSON body to the `tempData` struct
	if err := c.ShouldBindJSON(&tempData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
		INSERT INTO temperature_data (temperature, timestamp)
		VALUES ($1, $2)`


	_, err := db.Exec(sqlStatement, tempData.Temperature, tempData.Timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		log.Println("Failed to insert data:", err)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"temperature": tempData.Temperature,
		"timestamp":  tempData.Timestamp,
	})
}


func deleteTemperatureData(c *gin.Context) {
	id := c.Param("id")

	sqlStatement := `
		DELETE FROM temperature_data WHERE id=$1`


	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		log.Println("Failed to delete data:", err)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"id": id,
	})
}

func formatTimeStamp(timeToFormat string)(string, error) {
	fromTime, err := time.Parse(time.RFC3339, timeToFormat)
	if err !=nil {
		return "", err
	}

	return fromTime.Format("2006-01-02 15:04:05"), nil

}

func getTemperatureDataInSpan(c *gin.Context) {

	formattedFromTime, err := formatTimeStamp(c.Query("start"))
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp format"})
	}

	formattedToTime, err := formatTimeStamp(c.Query("end"))
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp format"})
	}

	sqlStatement := `
		SELECT temperature, timestamp FROM temperature_data WHERE "timestamp" BETWEEN $1 AND $2`

	result, err := db.Query(sqlStatement, formattedFromTime, formattedToTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		log.Println("Failed to get data:", err)
		return
	}

    // An temperatureData slice to hold data from returned rows.
    var temperatures []TemperatureData

    // Loop through rows, using Scan to assign column data to struct fields.
    for result.Next() {
        var temperature TemperatureData
        if err := result.Scan(&temperature.Temperature, &temperature.Timestamp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        temperatures = append(temperatures, temperature)
    }
    // Check for errors after the loop
    if err := result.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"result": temperatures,
	})

}
