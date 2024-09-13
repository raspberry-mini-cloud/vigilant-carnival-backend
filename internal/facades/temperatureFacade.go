package temperatureFacade

import (
	"log"
	"net/http"
	db "project/internal/db"
	. "project/internal/models"
	. "project/internal/utils"

	"github.com/gin-gonic/gin"
)



func GetTemperatureDataInSpan(c *gin.Context) {

	formattedFromTime, err := FormatTimeStamp(c.Query("start"))
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp format"})
	}

	formattedToTime, err := FormatTimeStamp(c.Query("end"))
	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp format"})
	}

	sqlStatement := `
		SELECT temperature, timestamp FROM temperature_data WHERE "timestamp" BETWEEN $1 AND $2`

	result, err := db.Db.Query(sqlStatement, formattedFromTime, formattedToTime)
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

func PostTemperatureData(c *gin.Context) {
	var tempData TemperatureData

	// Bind the JSON body to the `tempData` struct
	if err := c.ShouldBindJSON(&tempData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `
		INSERT INTO temperature_data (temperature, timestamp)
		VALUES ($1, $2)`


	_, err := db.Db.Exec(sqlStatement, tempData.Temperature, tempData.Timestamp)
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


func DeleteTemperatureData(c *gin.Context) {
	id := c.Param("id")

	sqlStatement := `
		DELETE FROM temperature_data WHERE id=$1`


	_, err := db.Db.Exec(sqlStatement, id)
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