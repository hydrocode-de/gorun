package frontend

import (
	"embed"
	"io/fs"
)

//go:embed manager/**/*
//go:embed manager/*
var manager embed.FS

func GetManager() fs.FS {
	managerFS, err := fs.Sub(manager, "manager")
	if err != nil {
		panic(err)
	}
	return managerFS
}
