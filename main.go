package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const (
		serverRoot = "."
		serverPort = ":8080"
	)
	var myEnv map[string]string

	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := sql.Open("postgres", myEnv["DB_URL"])
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	apiConf := apiConfig{db: database.New(db), config: myEnv}

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(serverRoot)))
	mux.Handle("/app/", apiConf.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/chirps", apiConf.chirpyHandler)
	mux.HandleFunc("POST /api/chirps", apiConf.validateChirpHandler)
	mux.HandleFunc("POST /api/users", apiConf.usersHandler)

	mux.HandleFunc("GET /admin/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.resetHandler)

	server := http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
