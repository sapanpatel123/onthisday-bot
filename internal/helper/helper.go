package helper

import (
	"os"
	"time"
)

// Exists returns true if a given path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}

// isDateSame returns true if the month and day are the same
func IsDateSame(reqDate time.Time, fileDate time.Time) bool {
	_, reqM, reqD := reqDate.Date()
	_, dateM, dateD := fileDate.Date()

	if (reqM == dateM) && (reqD == dateD) {
		return true
	}

	return false
}
