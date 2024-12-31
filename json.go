package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type error struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %d - %s", code, message)
		return
	}

	responsdWithJson(w, code, error{Error: message})
}

func responsdWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal error: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
