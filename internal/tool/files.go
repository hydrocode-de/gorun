package tool

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/hydrocode-de/gorun/internal/files"
)

func (t *Tool) ListResults() ([]files.ResultFile, error) {
	if t.Status != "finished" {
		return nil, errors.New("unfished tools cannot list results")
	}

	hostOut, ok := t.Mounts["/out"]
	if !ok {
		return nil, fmt.Errorf("tool %v did not mount /out. That means there is no folder with results", t.Name)
	}
	return files.ReadDir(hostOut, true, hostOut)
}

type WriteFileMeta struct {
	Filename string
	MimeType string
	FullPath string
}

func (t *Tool) WriteResultFile(resultPath string, w io.Writer) (*WriteFileMeta, error) {
	files, err := t.ListResults()
	if err != nil {
		return nil, err
	}

	var matchedFile string
	for _, file := range files {
		if strings.HasSuffix(resultPath, file.Name) {
			matchedFile = file.RelPath
			break
		}
	}
	if matchedFile == "" {
		return nil, fmt.Errorf("the result file %s was not found in the tool %s results", resultPath, t.Name)
	}

	fullPath := path.Join(t.Mounts["/out"], matchedFile)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	mimeType := http.DetectContentType(buffer)
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(w, file)
	if err != nil {
		return nil, err
	}

	return &WriteFileMeta{
		Filename: path.Base(matchedFile),
		MimeType: mimeType,
		FullPath: fullPath,
	}, nil
}
