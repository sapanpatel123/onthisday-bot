package search

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/sapanpatel123/onthisday-bot/internal/helper"
)

// FindPhotos gathers all files from the same month and day as today
func FindPhotos(path string, reqDate time.Time) ([]string, error) {
	var files []string
	var onthisdayFiles []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}

		if !info.IsDir() && (strings.ToLower(filepath.Ext(path)) == ".jpeg" || strings.ToLower(filepath.Ext(path)) == ".jpg") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return nil, err
	}

	for _, file := range files {
		t, err := times.Stat(file)
		if err != nil {
			log.Fatal(err.Error())
		}

		valid := helper.IsDateSame(reqDate, t.BirthTime())

		if valid {
			onthisdayFiles = append(onthisdayFiles, file)
		}
	}

	return onthisdayFiles, nil
}
