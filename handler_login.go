package main

import (
	"encoding/json"
	"net/http"

	"github.com/mohndakbar/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	responsdWithJson(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
