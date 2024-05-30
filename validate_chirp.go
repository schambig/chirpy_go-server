package main

import (
	"encoding/json"
	"net/http"
	// "log"
	"strings"
)

// struct for the json body to expect
type validChirp struct {
	Body string `json:"body"`
}

// struct to return marshaled JSON
type returnChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidChirp(w http.ResponseWriter, r *http.Request) {
	var chirp validChirp

	// decode the json request body into the chirp variable
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		// respond with error if json decoding fails
		respondWithError(w, http.StatusInternalServerError, "Something went wrong when decoding JSON")
		return
	}
	
	// check length of the chirp (Body field)
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		// respond with error if Body field exceeds length
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return		
	}

	newBody := replaceProfaneWords(w, chirp.Body)
	
	// respond with successful message if all went as expected
	// respondWithJSON(w, http.StatusOK, map[string]bool{"valid":true}) // a map can be marshalled
	respondWithJSON(w, http.StatusOK, returnChirp{
		CleanedBody: newBody,
	})
}

func replaceProfaneWords(w http.ResponseWriter, msg string) string {
	// log.Printf("   >>>THIS IS THE PAYLOAD BODY: %s", msg)
	strSplitted := strings.Split(msg, " ")
	// log.Printf("   >>>STR LOWERED AND SPLITTED: %s", strLowerAndSplitted)
	var newString []string
	for _, e := range strSplitted {
		if strings.ToLower(e) == "kerfuffle" || strings.ToLower(e) == "sharbert" || strings.ToLower(e) == "fornax" {
			newString = append(newString, "****")
		} else {
			newString = append(newString, e)
		}
	}
	// respondWithJSON(w, http.StatusOK, )
	// log.Println(strings.Join(newString, " "))
	return strings.Join(newString, " ")
}
