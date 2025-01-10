package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mohndakbar/chirpy/internal/auth"
	"github.com/mohndakbar/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type params struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params := params{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Failed to parse request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.dbQueires.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	responsdWithJson(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	err := cfg.dbQueires.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	responsdWithJson(w, 200, "All users deleted")
}
