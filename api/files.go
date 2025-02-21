package api

import (
	"net/http"

	"github.com/hydrocode-de/gorun/config"
)

// This function copies uploaded files into a temporary directory and returns the supplied file name and the path in a mapping
func HandleFileUpload(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {

}
