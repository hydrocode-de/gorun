package api

import (
	"encoding/json"
	"net/http"

	"github.com/hydrocode-de/gorun/config"
)

func HandleFuncWithConfig(handler func(http.ResponseWriter, *http.Request, *config.APIConfig), c *config.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

func CreateServer(c *config.APIConfig) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /runs", HandleFuncWithConfig(CreateRun, c))
	mux.HandleFunc("GET /runs/{id}", RunMiddleware(GetRunStatus, c))
	mux.HandleFunc("POST /runs/{id}/start", RunMiddleware(HandleRunStart, c))
	mux.HandleFunc("POST /files", HandleFuncWithConfig(HandleFileUpload, c))
	mux.HandleFunc("GET /specs", HandleFuncWithConfig(ListToolSpecs, c))
	mux.HandleFunc("GET /specs/{toolname}", HandleFuncWithConfig(GetToolSpec, c))

	return mux, nil
}

func RespondWithError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	w.Write([]byte(err))
}

func ResondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
