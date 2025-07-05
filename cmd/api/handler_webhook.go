package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Lyra-poing-serre/chirpy/internal/auth"
	"github.com/Lyra-poing-serre/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) WebhookHandler(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	_, err := auth.GetAPIKey(req.Header)
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
	if request.Event != "user.upgraded" {
		jsonResponse(w, http.StatusNoContent, "")
		return
	}
	err = a.Db.UpdateRedUser(context.Background(), database.UpdateRedUserParams{
		ID:          request.Data.UserId,
		IsChirpyRed: true,
	})
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	jsonResponse(w, http.StatusNoContent, "")
}
