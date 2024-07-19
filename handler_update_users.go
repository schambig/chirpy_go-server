package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"

	"github.com/schambig/chirpy_go-server/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

/* 	// from docs: https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Token
	type Token struct {
		Raw       string                 // Raw contains the raw token.  Populated when you [Parse] a token
		Method    SigningMethod          // Method is the signing method used or to be used
		Header    map[string]interface{} // Header is the first segment of the token in decoded form
		Claims    Claims                 // Claims is the second segment of the token in decoded form
		Signature []byte                 // Signature is the third segment of the token in decoded form.  Populated when you Parse a token
		Valid     bool                   // Valid specifies if the token is valid.  Populated when you Parse/Verify a token
	} */

	decoder := 	json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
		return			
	}

	// header format: Authorization: Bearer <token>, we just need the token
	// If there are no values associated with the key, Get returns ""
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "No Authorization header provided")
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format in AuthZ header")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

	// validate the signature of the JWT and extract the claims
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid or has expired")
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	// updating the user
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid user ID in token")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.UpdateUser(userID, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return		
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
	})
}
