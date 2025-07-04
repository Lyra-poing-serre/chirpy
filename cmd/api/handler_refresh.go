package api

import (
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
)

func (a *ApiConfig) RefreshHandler(w http.ResponseWriter, req *http.Request) {
	type refreshResponse struct {
		Token string `json:"token"`
	}
	uID, err := a.verifyRefreshToken(req.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	accessToken, err := auth.MakeJWT(uID, a.Config["SERVER_SECRET"], time.Hour*60)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, refreshResponse{Token: accessToken})
}
