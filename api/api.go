package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/frontend"
)

func HandleApiKey(handler func(http.ResponseWriter, *http.Request), c *config.APIConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		if apiKey != "" {
			userId, err := auth.ValidateJWT(apiKey, c.Secret)
			if err == nil {
				r.Header.Set("X-User-ID", userId)
			}
		}

		handler(w, r)
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

	mux.HandleFunc("GET /runs", HandleApiKey(HandleFuncWithConfig(GetAllRuns, c), c))
	mux.HandleFunc("POST /runs", HandleApiKey(HandleFuncWithConfig(CreateRun, c), c))
	mux.HandleFunc("GET /runs/{id}", HandleApiKey(RunMiddleware(GetRunStatus, c), c))
	mux.HandleFunc("DELETE /runs/{id}", HandleApiKey(RunMiddleware(DeleteRun, c), c))
	mux.HandleFunc("POST /runs/{id}/start", HandleApiKey(RunMiddleware(HandleRunStart, c), c))
	mux.HandleFunc("GET /runs/{id}/results", HandleApiKey(RunMiddleware(ListRunResults, c), c))
	mux.HandleFunc("GET /runs/{id}/results/{filename}", HandleApiKey(RunMiddleware(GetResultFile, c), c))
	mux.HandleFunc("POST /files", HandleApiKey(HandleFuncWithConfig(HandleFileUpload, c), c))
	mux.HandleFunc("GET /files", HandleApiKey(HandleFuncWithConfig(FindFile, c), c))
	mux.HandleFunc("GET /specs", HandleFuncWithConfig(ListToolSpecs, c))
	mux.HandleFunc("GET /specs/{toolname}", HandleFuncWithConfig(GetToolSpec, c))
	mux.HandleFunc("POST /auth/refresh", HandleFuncWithConfig(HandleRefreshToken, c))
	mux.HandleFunc("POST /auth/login", HandleFuncWithConfig(HandleLogin, c))
	return mux, nil
}

func RespondWithError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(err))
}

func ResondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}
