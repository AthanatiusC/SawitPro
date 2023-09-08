package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Isolate token from string, returns valid token out of auth bearer
func getToken(auth string) (string, error) {
	// Bearer {{$token}}, split into 2 index, get second index in this case 1st index
	jwtToken := strings.Split(auth, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("invalid token")
	}

	return jwtToken[1], nil
}

// Validate JWT using envar secrets, returns valid jwt token
func (s *Server) ValidateJWT(authParam string) (token *jwt.Token) {
	auth, err := getToken(authParam)
	if err != nil {
		return
	}

	token, err = jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) { // JWT Parse require func(interface{},error) as its argument
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return
	}

	return
}

// Get JWT claim value by string key, returns string and error
func (s *Server) GetJWTClaims(token *jwt.Token, key string) (string, error) {
	claims := token.Claims.(jwt.MapClaims)[key].(string)
	return claims, nil
}

// Generate JWT using envar secrets, returns valid jwt token and error
func (s *Server) GenerateJWT(id int) (token string, err error) {
	exp := time.Now().Add(time.Hour * 24).Unix() // Common exp time
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  fmt.Sprint(id), // to ensure string convert when get claims
		"exp": exp,
	})

	token, err = claims.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return
	}

	return
}
