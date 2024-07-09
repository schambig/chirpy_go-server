package main

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/schambig/chirpy_go-server/internal/auth"
	"github.com/schambig/chirpy_go-server/internal/database"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
		return		
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "Email already in use")
			return	
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return		
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
	})
}
