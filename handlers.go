package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, req *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(h)
}

func plainTextHandler(data string) http.Handler {
	if data == "" {
		log.Fatalln("Empty data to write")
	}
	f := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, data)
	}
	return http.HandlerFunc(f)
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	h := plainTextHandler("OK")
	h.ServeHTTP(w, req)
}

func (a *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	h := plainTextHandler(fmt.Sprintf("Hits: %d", a.fileserverHits.Load()))
	h.ServeHTTP(w, req)
}

func (a *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	a.fileserverHits.Store(0)
	h := plainTextHandler("Hits are reseted")
	h.ServeHTTP(w, req)
}
