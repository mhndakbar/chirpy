package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueires.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	jsonChirps := []Chirp{}
	for _, chirp := range chirps {
		jsonChirps = append(jsonChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	responsdWithJson(w, 200, jsonChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, 400, "Invalid chirp id")
		return
	}

	chirp, err := cfg.dbQueires.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if chirp.ID == uuid.Nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	responsdWithJson(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
