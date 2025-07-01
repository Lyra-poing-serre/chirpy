package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

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
		Body string `json:"body"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	var err error

	if err = decoder.Decode(&params); err != nil {
		jsonResponse(w, http.StatusInternalServerError, err)
	} else if len(params.Body) > 140 {
		jsonResponse(w, http.StatusInternalServerError, errors.New("chirp is too long"))
	} else if params.Body == "" {
		jsonResponse(w, http.StatusBadRequest, errors.New("empty body"))
	} else {
		params.Body = cleanChirp(params.Body)
		jsonResponse(w, 200, params)
	}

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

func (a *apiConfig) validateUsersHandler(w http.ResponseWriter, req *http.Request) {
	if a.config["PLATFORM"] != "dev" {
		errorResponse(w, http.StatusForbidden, "PLATFORM != dev")
		return
	}
	type reqParameters struct {
		Email string `json:"email"`
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

	dbUsr, err := a.db.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     data.Email,
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
