package utils

import (
	"fmt"
	"time"
)

// iki gae generate nomer antrian
func GenerateNoAntrian(count int, weekday time.Weekday) string {
	prefix := "A"
	switch weekday {
	case time.Sunday:
		prefix = "M"
	case time.Monday:
		prefix = "S"
	case time.Tuesday:
		prefix = "SL"
	case time.Wednesday:
		prefix = "R"
	case time.Thursday:
		prefix = "K"
	case time.Friday:
		prefix = "J"
	case time.Saturday:
		prefix = "ST"
	}
	return fmt.Sprintf("%s-%02d", prefix, count+1)
}
