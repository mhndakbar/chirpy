package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type validateChirpRequest struct {
	Body string `json:"body"`
}

type cleanedBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	chirp := validateChirpRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 500, "Failed to parse request body")
		return
	}

	const chirpMaxLength = 140
	if len(chirp.Body) > chirpMaxLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	msg := cleanBody(badWords, chirp.Body)
	responsdWithJson(w, 200, cleanedBody{CleanedBody: msg})
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
