package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		serverRoot = "."
		serverPort = ":8080"
	)
	mux := http.NewServeMux()
	apiConf := apiConfig{}

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(serverRoot)))
	mux.Handle("/app/", apiConf.middlewareMetricsInc(fileHandler))
	mux.HandleFunc("GET /admin/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.resetHandler)

	server := http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
