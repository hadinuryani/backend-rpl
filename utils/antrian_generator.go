package utils

import "fmt"

// GenerateNoAntrian creates a queue number in the format "A-01", "A-02", etc.
// count is the number of existing antrian for today (before this new one).
func GenerateNoAntrian(count int) string {
	return fmt.Sprintf("A-%02d", count+1)
}
