package temperatureData

import "time"

type TemperatureData struct {
	Temperature float64 `json:"temperature"`
	Timestamp time.Time  `json:"timestamp"`
	//Location string `json:"location"`
}