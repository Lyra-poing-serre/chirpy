package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) UsersHandler(w http.ResponseWriter, req *http.Request) {
	if a.Config["PLATFORM"] != "dev" {
		errorResponse(w, http.StatusForbidden, "PLATFORM != dev")
		return
	}
	type reqParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type jsonUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	var data reqParameters

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&data)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	hash, err := auth.HashPassword(data.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	dbUsr, err := a.Db.CreateUser(req.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		HashedPassword: hash,
		Email:          data.Email,
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusCreated, jsonUser{
		ID:        dbUsr.ID,
		CreatedAt: dbUsr.CreatedAt,
		UpdatedAt: dbUsr.UpdatedAt,
		Email:     dbUsr.Email,
	})
}
