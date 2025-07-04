package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
)

func (a *ApiConfig) RevokeHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	_, err = a.verifyRefreshToken(req.Header)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = a.Db.RevokeRefreshToken(context.Background(), database.RevokeRefreshTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusNoContent, "")
}
