package main

import (
	"net/http"
	"time"

	"github.com/mohndakbar/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	refreshToken, err := cfg.dbQueires.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	userID, err := cfg.dbQueires.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	jwtToken, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Duration(60)*time.Minute)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	responsdWithJson(w, http.StatusOK, map[string]string{
		"token": jwtToken,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToekn(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	_, err = cfg.dbQueires.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can't validate token")
		return
	}

	responsdWithJson(w, http.StatusNoContent, nil)
}
