package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

func (a *ApiConfig) ChirpyHandler(w http.ResponseWriter, req *http.Request) {
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
	chirp, err := a.Db.GetChirps(context.Background(), cID)
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
