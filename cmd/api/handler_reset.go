package api

import (
	"context"
	"net/http"
)

func (a *ApiConfig) ResetHandler(w http.ResponseWriter, req *http.Request) {
	err := a.Db.ResetUsers(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.FileserverHits.Store(0)
	h := PlainTextHandler("text/plain; charset=utf-8", "Hits are reseted\nUsers table is now empty.")
	h.ServeHTTP(w, req)
}
