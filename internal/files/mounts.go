package files

import (
	"os"
	"path"

	"github.com/hydrocode-de/gorun/internal/helper"
)

func CreateNewMountPaths(mountPath string, level string) map[string]string {
	mounts := make(map[string]string)

	if level == "_random" {
		level = helper.GetRandomString(12)
	}
	mounts["/in"] = path.Join(mountPath, level, "in")
	mounts["/out"] = path.Join(mountPath, level, "out")

	for _, hostPath := range mounts {
		os.MkdirAll(hostPath, 0755)
	}

	return mounts
}
