package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/frontend"
)

func HandleRequireApiKey(handler func(http.ResponseWriter, *http.Request) error, c *config.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		if authHeader == "" || apiKey == "" {
			RespondWithError(w, http.StatusUnauthorized, "missing api key in Authorization header")
			return
		}

		// the validation function is returning an error if the key is invalid
		// sorry this look maybe a bit confusing compared to other go code
		err := auth.ValidateApiKey(apiKey, r.Context(), c.GetDB())
		if err != nil {
			handler(w, r)
		} else {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
		}

	}
}

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

	// add a FileServer to serve the manager
	//mux.Handle("/manager/", http.StripPrefix("/manager/", http.FileServer(http.Dir("manager/build"))))
	mux.Handle("/manager/", http.StripPrefix("/manager/", http.FileServerFS(frontend.GetManager())))

	mux.HandleFunc("GET /runs", HandleFuncWithConfig(GetAllRuns, c))
	mux.HandleFunc("POST /runs", HandleFuncWithConfig(CreateRun, c))
	mux.HandleFunc("GET /runs/{id}", RunMiddleware(GetRunStatus, c))
	mux.HandleFunc("DELETE /runs/{id}", RunMiddleware(DeleteRun, c))
	mux.HandleFunc("POST /runs/{id}/start", RunMiddleware(HandleRunStart, c))
	mux.HandleFunc("GET /runs/{id}/results", RunMiddleware(ListRunResults, c))
	mux.HandleFunc("GET /runs/{id}/results/{filename}", RunMiddleware(GetResultFile, c))
	mux.HandleFunc("POST /files", HandleFuncWithConfig(HandleFileUpload, c))
	mux.HandleFunc("GET /files", HandleFuncWithConfig(FindFile, c))
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
