package files

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/hydrocode-de/gorun/config"
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
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func Cleanup(c *config.APIConfig) error {
	uploadDir := path.Join(c.BaseTempDir, "uploads")
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return err
	}
	err = filepath.WalkDir(uploadDir, func(p string, d fs.DirEntry, e error) error {
		return walkDirCheckTimestamp(p, d, e, time.Now().Add(-c.MaxTempAge))
	})
	if err != nil {
		return err
	}

	return nil
}
