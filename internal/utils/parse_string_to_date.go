package utils

import "time"

func ParseStringToDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}

	parsedTime, err := time.Parse(time.RFC3339, *dateStr)
	if err != nil {
		return nil
	}
	return &parsedTime
}
