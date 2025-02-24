package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/hydrocode-de/gorun/config"
)

// This function copies uploaded files into a temporary directory and returns the supplied file name and the path in a mapping
func HandleFileUpload(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	if err := r.ParseMultipartForm(c.MaxUploadSize); err != nil {
		RespondWithError(w, 413, fmt.Sprintf("error parsing multipart form: %s", err))
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error reading uploaded file: %s", err))
	}
	defer file.Close()

	tempBaseDir := path.Join(c.BaseTempDir, "uploads")
	err = os.MkdirAll(tempBaseDir, 0755)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("error creating gorun temporary directory base: %s", err))
	}
	tempDir, err := os.MkdirTemp(tempBaseDir, "")
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("error creating temporary directory: %s", err))
	}
	targetPath := path.Join(tempDir, handler.Filename)
	openf, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("error creating target file: %s", err))
	}
	defer openf.Close()
	writtenBytes, err := io.Copy(openf, file)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("error writing to target file: %s", err))
	}

	ResondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"path": targetPath,
		"size": writtenBytes,
		"name": handler.Filename,
		"type": handler.Header.Get("Content-Type"),
	})
}
