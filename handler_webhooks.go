package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg apiConfig) handlerCreateWebhook(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserId uuid.UUID `json:"user_id"`
	}
	type params struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	reqParams := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse request body")
		return
	}

	if reqParams.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueires.UpgradeUserToChirpyRed(r.Context(), reqParams.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	responsdWithJson(w, http.StatusNoContent, nil)
}
