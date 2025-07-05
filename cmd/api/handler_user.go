package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

type reqParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type jsonUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (a *ApiConfig) UsersHandler(w http.ResponseWriter, req *http.Request) {
	if a.Config["PLATFORM"] != "dev" {
		errorResponse(w, http.StatusForbidden, "PLATFORM != dev")
		return
	}
	switch req.Method {
	case "POST":
		a.CreateUserHandler(w, req)
	case "PUT":
		a.UpdateUserHandler(w, req)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}

func (a *ApiConfig) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
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
		ID:          dbUsr.ID,
		CreatedAt:   dbUsr.CreatedAt,
		UpdatedAt:   dbUsr.UpdatedAt,
		Email:       dbUsr.Email,
		IsChirpyRed: dbUsr.IsChirpyRed,
	})
}

func (a *ApiConfig) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	uId, err := auth.ValidateJWT(token, a.Config["SERVER_SECRET"])
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request reqParameters
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashPwd, err := auth.HashPassword(request.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	dbUsr, err := a.Db.UpdateUserPwdEmail(context.Background(), database.UpdateUserPwdEmailParams{
		ID:             uId,
		Email:          request.Email,
		HashedPassword: hashPwd,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, jsonUser{
		ID:          dbUsr.ID,
		CreatedAt:   dbUsr.CreatedAt,
		UpdatedAt:   dbUsr.UpdatedAt,
		Email:       dbUsr.Email,
		IsChirpyRed: dbUsr.IsChirpyRed,
	})
}
