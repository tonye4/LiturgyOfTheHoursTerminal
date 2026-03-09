package main

import (
	"strings"
	"time"
)

// Date format within the JSON file: 20260206
func today() string {
	currentDate := time.Now().Format(time.DateOnly)
	currentDateFormatted := strings.ReplaceAll(currentDate, "-", "")

	return currentDateFormatted
}
