package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {
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
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	retrievedId, err := auth.ValidateJWT(token, a.Config["SERVER_SECRET"])
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}

	if err = decoder.Decode(&params); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(params.Body) > 140 {
		jsonResponse(w, http.StatusInternalServerError, errors.New("chirp is too long"))
		return
	} else if params.Body == "" {
		jsonResponse(w, http.StatusBadRequest, errors.New("empty body"))
		return
	}
	chirp, err := a.Db.CreateChirp(context.Background(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanChirp(params.Body),
		UserID:    retrievedId,
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
