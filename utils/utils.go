package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"jwt-course/models"
	"log"
	"net/http"
	"os"
	"strings"
)

func RespondWithError(rw http.ResponseWriter, status int, responseError models.Error) {
	logger := log.New(os.Stdout, "writing data: ", log.LstdFlags)
	rw.WriteHeader(status)
	err := json.NewEncoder(rw).Encode(responseError)
	if err != nil {
		logger.Fatalln(err)
	}
}

func GenerateToken(user models.User) (string, error) {
	secret := os.Getenv("SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss": "course",
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var responseError models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				responseError.Message = err.Error()
				RespondWithError(rw, http.StatusUnauthorized, responseError)
				return
			}

			if token.Valid {
				next.ServeHTTP(rw, r)
			} else {
				responseError.Message = "invalid token"
				RespondWithError(rw, http.StatusUnauthorized, responseError)
				return
			}
		} else {
			responseError.Message = "invalid token"
			RespondWithError(rw, http.StatusUnauthorized, responseError)
			return
		}
	}
}