package main

import (
	"strings"
	"time"
)

// Date format within the JSON file: 20260206
func Today() string {
	currentDate := time.Now().Format(time.DateOnly)
	currentDateFormatted := strings.ReplaceAll(currentDate, "-", "")

	return currentDateFormatted
}

// create helper for viewing specific days in the prayer cache.
// Make it indented.
