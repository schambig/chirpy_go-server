package auth

import (
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// compares a bcrypt hashed password with its possible plaintext equivalent,
// returns nil on success, or an error on failure
func CheckHashPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// create JWT using JWT library: https://github.com/golang-jwt/jwt 
func MakeJWT(userID int, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: fmt.Sprintf("%d", userID),
	})

	return token.SignedString(signingKey)
}
