package api

import (
	"io"
	"log"
	"net/http"
)

func PlainTextHandler(contentType string, data string) http.Handler {
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
