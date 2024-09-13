package utils

import "time"

func FormatTimeStamp(timeToFormat string)(string, error) {
	fromTime, err := time.Parse(time.RFC3339, timeToFormat)
	if err !=nil {
		return "", err
	}

	return fromTime.Format("2006-01-02 15:04:05"), nil

}