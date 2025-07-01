package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/Lyra-poing-serre/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	config         map[string]string
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, req *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(h)
}

func plainTextHandler(contentType string, data string) http.Handler {
	if data == "" {
		log.Fatalln("Empty data to write")
	}
	f := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", contentType) // "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, data)
	}
	return http.HandlerFunc(f)
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	h := plainTextHandler("text/plain; charset=utf-8", "OK\n")
	h.ServeHTTP(w, req)
}

func (a *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	h := plainTextHandler("text/html", fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, a.fileserverHits.Load()))
	h.ServeHTTP(w, req)
}

func (a *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	err := a.db.ResetUsers(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
	}
	a.fileserverHits.Store(0)
	h := plainTextHandler("text/plain; charset=utf-8", "Hits are reseted\nUsers table is now empty.")
	h.ServeHTTP(w, req)
}
