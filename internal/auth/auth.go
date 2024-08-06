package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthzHeaderIncluded = errors.New("no authz header included in request")

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

// extract the token from authorization header
func GetBearerToken(headers http.Header) (string, error) {
	/* // alternative to get the token string
	authzHeader := headers.Get("Authorization")
	if authzHeader == "" {
		return "", errors.New("no Authorization header provided")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authzHeader, bearerPrefix) {
		return "", errors.New("invalid token format in authorization header")
	}
	tokenStr := strings.TrimPrefix(authzHeader, bearerPrefix)
	return tokenStr, nil */

	authzHeader := headers.Get("Authorization")
	if authzHeader == "" {
		return "", ErrNoAuthzHeaderIncluded
	}

	splitAuthz := strings.Split(authzHeader, " ")
	if len(splitAuthz) < 2 || splitAuthz[1] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuthz[1], nil
}

// validate the JWT, use the jwt.ParseWithClaims function to validate the signature of the JWT
// and extract the claims into a *jwt.Token struct
// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#ParseWithClaims
func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(t *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string("chirpy") {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}
