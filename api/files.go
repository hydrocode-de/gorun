package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/hydrocode-de/gorun/internal/files"
	"github.com/spf13/viper"
)

type FindFilesResponse struct {
	Count int                `json:"count"`
	Files []files.ResultFile `json:"files"`
}

// This function copies uploaded files into a temporary directory and returns the supplied file name and the path in a mapping
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	maxUploadSize := viper.GetInt("max_upload_size")
	tempPath := viper.GetString("temp_path")

	if err := r.ParseMultipartForm(int64(maxUploadSize)); err != nil {
		RespondWithError(w, 413, fmt.Sprintf("error parsing multipart form: %s", err))
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error reading uploaded file: %s", err))
	}
	defer file.Close()

	tempBaseDir := path.Join(tempPath, "uploads")
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

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"path": targetPath,
		"size": writtenBytes,
		"name": handler.Filename,
		"type": handler.Header.Get("Content-Type"),
	})
}

func FindFile(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		RespondWithError(w, http.StatusBadRequest, "missing pattern, you need to provide a 'pattern' query parameter")
		return
	}

	target := r.URL.Query().Get("target")
	if target == "" {
		target = "both"
	}

	if filepath.Ext(pattern) == "" {
		pattern += "*.*"
	}

	mountPath := viper.GetString("mount_path")
	matches, err := files.Find(pattern, mountPath, files.Target(target))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, FindFilesResponse{
		Count: len(matches),
		Files: matches,
	})

}
