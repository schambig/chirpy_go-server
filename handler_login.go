package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/schambig/chirpy_go-server/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
		Token string `json:"token,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
		return		
	}
	
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckHashPassword(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password didn't match, try again")
		return
	}

	// if ommited in the JSON, ExpiresInSeconds will be defaulted to 0 (omitempty directive)
	expiresInSecondsDefault := 24 * 60 * 60 // 86,400 seconds (1 day)
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = expiresInSecondsDefault
	} else if params.ExpiresInSeconds > expiresInSecondsDefault {
		params.ExpiresInSeconds = expiresInSecondsDefault
	}

	token, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}
