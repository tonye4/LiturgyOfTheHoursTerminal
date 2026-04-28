package main

import (
	"strings"
	"time"

	"github.com/tonye4/LiturgyOfTheHoursTerminal/tabs"
)

// Date format within the JSON file: 20260206
func GetFormattedDate(dayOpt tabs.Tab) string {

	switch dayOpt {
	case tabs.Today:
		currentDate := time.Now().Format(time.DateOnly)
		currentDateFormatted := strings.ReplaceAll(currentDate, "-", "")

		return currentDateFormatted

	case tabs.Yesterday:
		yesterday := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
		yesterdayDateFormatted := strings.ReplaceAll(yesterday, "-", "")

		return yesterdayDateFormatted

	case tabs.Tomorrow:
		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)
		tomorrowDateFormatted := strings.ReplaceAll(tomorrow, "-", "")

		return tomorrowDateFormatted
	default:
		errorString := ""
		return errorString
	}
}

// create helper for viewing specific days in the prayer cache.
// Make it indented.
