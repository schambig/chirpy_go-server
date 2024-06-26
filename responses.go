package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Server error 5XX: %s,", msg)
	} else if code >= 400 {
		log.Printf("Client error 4XX: %s,", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
	
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		log.Printf("Error while marshalling JSON: %s,", err)
		w.WriteHeader(500) // marshalling error is a server error
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
