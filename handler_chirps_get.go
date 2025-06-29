package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/mohndakbar/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort")

	userID := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if userID != "" {
		userID, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(w, 400, "Invalid author_id")
			return
		}
		chirps, err = cfg.dbQueires.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	} else {
		chirps, err = cfg.dbQueires.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	}

	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
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
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	chirp, err := cfg.dbQueires.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
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
