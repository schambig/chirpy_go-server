package main

import (
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerGetChirpID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
        respondWithError(w, http.StatusBadRequest, "Invalid request path")
        return
	}
	
	chirpIDStr := pathParts[3]
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
