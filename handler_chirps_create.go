package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/mohndakbar/chirpy/internal/auth"
	"github.com/mohndakbar/chirpy/internal/database"
)

const chirpMaxLength = 140

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
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

	type chirpRequest struct {
		Body string `json:"body"`
	}

	chirpParams := chirpRequest{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirpParams)
	if err != nil {
		respondWithError(w, 500, "Failed to parse request body")
		return
	}

	cleaned, err := validate(chirpParams.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	NewChirp, err := cfg.dbQueires.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	responsdWithJson(w, 201, Chirp{
		ID:        NewChirp.ID,
		CreatedAt: NewChirp.CreatedAt,
		UpdatedAt: NewChirp.UpdatedAt,
		Body:      NewChirp.Body,
		UserID:    NewChirp.UserID,
	})
}

func validate(msg string) (string, error) {
	if len(msg) > chirpMaxLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	msg = cleanBody(badWords, msg)

	return msg, nil
}

func cleanBody(words map[string]struct{}, body string) string {
	msgSlice := strings.Split(body, " ")
	for i, word := range msgSlice {
		wordLowered := strings.ToLower(word)
		if _, ok := words[wordLowered]; ok {
			msgSlice[i] = "****"
		}
	}
	cleanedMsg := strings.Join(msgSlice, " ")
	return cleanedMsg
}
