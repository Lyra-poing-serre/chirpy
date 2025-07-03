package api

import "net/http"

func ReadinessHandler(w http.ResponseWriter, req *http.Request) {
	h := PlainTextHandler("text/plain; charset=utf-8", "OK\n")
	h.ServeHTTP(w, req)
}
