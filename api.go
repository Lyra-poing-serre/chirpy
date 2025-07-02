package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

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

func (a *apiConfig) validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}
	type jsonChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	var err error

	if err = decoder.Decode(&params); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	uId, err := uuid.Parse(params.UserId)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = a.db.GetUser(context.Background(), uId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, fmt.Errorf("%s is unknown", uId))
		return
	}
	if len(params.Body) > 140 {
		jsonResponse(w, http.StatusInternalServerError, errors.New("chirp is too long"))
		return
	} else if params.Body == "" {
		jsonResponse(w, http.StatusBadRequest, errors.New("empty body"))
		return
	}
	chirp, err := a.db.CreateChirp(context.Background(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanChirp(params.Body),
		UserID:    uId,
	})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusCreated, jsonChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func cleanChirp(body string) string {
	text := []string{}

	for word := range strings.SplitSeq(string(body), " ") {
		var match bool
		for _, banWord := range [3]string{"kerfuffle", "sharbert", "fornax"} {
			if strings.ToLower(word) == banWord {
				match = true
			}
		}
		if match {
			text = append(text, "****")
		} else {
			text = append(text, word)
		}
	}

	return strings.Join(text, " ")
}

func (a *apiConfig) usersHandler(w http.ResponseWriter, req *http.Request) {
	if a.config["PLATFORM"] != "dev" {
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
		Password  string    `json:"password"`
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
	dbUsr, err := a.db.CreateUser(req.Context(), database.CreateUserParams{
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

func (a *apiConfig) chirpyHandler(w http.ResponseWriter, req *http.Request) {
	cId := req.PathValue("chirpID")
	if cId == "" {
		errorResponse(w, http.StatusNotFound, "Not found")
	}
	type jsonChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	cID, err := uuid.Parse(cId)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := a.db.GetChirps(context.Background(), cID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, jsonChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (a *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type jsonUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Password  string    `json:"password"`
		Email     string    `json:"email"`
	}
	var request reqParameters

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := a.db.GetUserByEmail(context.Background(), request.Email)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if auth.CheckPasswordHash(request.Password, dbUser.HashedPassword) != nil {
		jsonResponse(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	jsonResponse(w, http.StatusOK, jsonUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
