package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Config         map[string]string
}

func (a *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, req *http.Request) {
		a.FileserverHits.Add(1)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(h)
}

func errorResponse(w http.ResponseWriter, statusCode int, err string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	jsonResponse(w, statusCode, errorResponse{Error: err})
}

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {

	out, err := json.Marshal(payload)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}

func (a *ApiConfig) verifyRefreshToken(head http.Header) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(head)
	if err != nil {
		return uuid.UUID{}, err
	}
	dbToken, err := a.Db.GetRefreshToken(context.Background(), token)
	if err != nil {
		return uuid.UUID{}, err
	}
	if dbToken.ExpiredAt.Before(time.Now()) || dbToken.RevokedAt.Valid {
		return uuid.UUID{}, errors.New("expired or revoked token")
	}
	return dbToken.UserID, nil
}
