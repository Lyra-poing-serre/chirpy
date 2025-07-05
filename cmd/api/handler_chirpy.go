package api

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

type jsonChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (a *ApiConfig) ChirpsHandler(w http.ResponseWriter, req *http.Request) {
	uId := req.URL.Query().Get("author_id")
	var chirpList []jsonChirp
	var chirpsDb []database.Chirp
	var err error
	if uId == "" {
		chirpsDb, err = a.Db.GetChirps(context.Background())
		if err != nil {
			errorResponse(w, http.StatusNotFound, err.Error())
			return
		}
	} else {
		userId, err := uuid.Parse(uId)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		chirpsDb, err = a.Db.GetChirpByAuthor(context.Background(), userId)
		if err != nil {
			errorResponse(w, http.StatusNotFound, err.Error())
			return
		}
	}
	sortingParam := req.URL.Query().Get("sort")
	for _, chirp := range chirpsDb {
		chirpList = append(chirpList, jsonChirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	if sortingParam == "desc" {
		sort.Slice(chirpList, func(a, b int) bool {
			return chirpList[a].CreatedAt.After(chirpList[b].CreatedAt)
		})
	}
	jsonResponse(w, http.StatusOK, chirpList)
}

func (a *ApiConfig) ChirpyHandler(w http.ResponseWriter, req *http.Request) {
	cID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	chirp, err := a.Db.GetChirpById(context.Background(), cID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
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

func (a *ApiConfig) RemoveChirpyHandler(w http.ResponseWriter, req *http.Request) {
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

	cID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	chirp, err := a.Db.GetChirpById(context.Background(), cID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	} else if uId != chirp.UserID {
		errorResponse(w, http.StatusForbidden, "not the current user")
		return
	}

	err = a.Db.DeleteChirp(context.Background(), chirp.ID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse(w, http.StatusNoContent, "")
}
