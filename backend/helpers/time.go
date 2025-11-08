package helpers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewDateTime converts time.Time to primitive.DateTime
func NewDateTime(t time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(t)
}

// NowDateTime returns current time as primitive.DateTime
func NowDateTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

// ParseDate parses YYYY-MM-DD format to time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate formats time.Time to YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseDateTime parses YYYY-MM-DD HH:MM:SS format to time.Time
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dateTimeStr)
}

// FormatDateTime formats time.Time to YYYY-MM-DD HH:MM:SS
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// CalculateWorkHours calculates hours between two time strings (HH:MM:SS)
func CalculateWorkHours(checkIn, checkOut string) (float64, error) {
	layout := "15:04:05"

	inTime, err := time.Parse(layout, checkIn)
	if err != nil {
		return 0, err
	}

	outTime, err := time.Parse(layout, checkOut)
	if err != nil {
		return 0, err
	}

	duration := outTime.Sub(inTime)
	hours := duration.Hours()

	return hours, nil
}
