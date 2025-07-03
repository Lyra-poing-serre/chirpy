package api

import (
	"fmt"
	"net/http"
)

func (a *ApiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	h := PlainTextHandler("text/html", fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, a.FileserverHits.Load()))
	h.ServeHTTP(w, req)
}
