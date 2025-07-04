package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Lyra-poing-serre/chirpy/cmd/api"
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
	apiConf := api.ApiConfig{Db: database.New(db), Config: myEnv}

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(serverRoot)))
	mux.Handle("/app/", apiConf.MiddlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.ChirpyHandler)
	mux.HandleFunc("POST /api/chirps", apiConf.ValidateChirpHandler)
	mux.HandleFunc("/api/users", apiConf.UsersHandler)
	mux.HandleFunc("POST /api/login", apiConf.LoginHandler)
	mux.HandleFunc("POST /api/refresh", apiConf.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiConf.RevokeHandler)

	mux.HandleFunc("GET /admin/healthz", api.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.ResetHandler)

	server := http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
