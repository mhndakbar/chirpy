package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mohndakbar/chirpy/internal/auth"
)

const (
	DefaultExpiresInSeconds = 3600
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	reqParams := params{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, 500, "Failed to parse request body")
		return
	}

	user, err := cfg.dbQueires.GetUserByEmail(r.Context(), reqParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(reqParams.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if reqParams.ExpiresInSeconds == 0 || reqParams.ExpiresInSeconds > DefaultExpiresInSeconds {
		reqParams.ExpiresInSeconds = DefaultExpiresInSeconds
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(reqParams.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	responsdWithJson(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwtToken,
	})

}
