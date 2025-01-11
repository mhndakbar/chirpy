package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mohndakbar/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate JWT token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate JWT token")
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	chirp, err := cfg.dbQueires.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	err = cfg.dbQueires.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	responsdWithJson(w, http.StatusNoContent, nil)
}
