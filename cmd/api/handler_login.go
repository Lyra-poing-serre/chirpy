package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (a *ApiConfig) LoginHandler(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Password        string `json:"password"`
		Email           string `json:"email"`
		ExpireInSeconds int    `json:"expires_in_seconds"`
	}
	type jsonUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	var request reqParameters

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := a.Db.GetUserByEmail(context.Background(), request.Email)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if auth.CheckPasswordHash(request.Password, dbUser.HashedPassword) != nil {
		jsonResponse(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	tokenLifeTime := time.Hour
	if request.ExpireInSeconds > 0 && time.Duration(request.ExpireInSeconds).Seconds() < tokenLifeTime.Seconds() {
		tokenLifeTime = time.Duration(request.ExpireInSeconds)
	}
	token, err := auth.MakeJWT(dbUser.ID, a.Config["SERVER_SECRET"], tokenLifeTime)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, jsonUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	})
}
