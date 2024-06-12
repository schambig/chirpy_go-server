package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"errors"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	defer r.Boby.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
		return		
	}

	cleaned, err := validateChirp(params.Boby)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return		
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return		
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		Boby: chirp.Boby,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	cleaned := replaceProfaneWords(body, badWords)
	return cleaned, nil	
}

func replaceProfaneWords(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
