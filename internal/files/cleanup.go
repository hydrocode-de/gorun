package files

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func walkDirCheckTimestamp(path string, d fs.DirEntry, err error, maxAge time.Time) error {
	if err != nil {
		return err
	}

	if !d.IsDir() {
		return nil
	}

	if filepath.Base(path) == "uploads" {
		return nil
	}

	// this is a drectory we need to check
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	shouldRemove := false

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(maxAge) {
			shouldRemove = true
			break
		}
	}

	if shouldRemove {
		fmt.Printf("Removing %v now....\n", path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		return filepath.SkipDir
	}

	return nil
}

func Cleanup() error {
	baseDir := viper.GetString("path")
	maxAge := viper.GetDuration("max_temp_age")

	uploadDir := path.Join(baseDir, "uploads")
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return err
	}
	err = filepath.WalkDir(uploadDir, func(p string, d fs.DirEntry, e error) error {
		return walkDirCheckTimestamp(p, d, e, time.Now().Add(maxAge))
	})
	if err != nil {
		return err
	}

	return nil
}
