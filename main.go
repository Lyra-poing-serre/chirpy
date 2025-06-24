package main

import (
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.NewServeMux(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
