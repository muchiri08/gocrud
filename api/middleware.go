package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"net/http"
)

func jWTMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the token from request
		tokenString, err := request.HeaderExtractor{"access-token"}.ExtractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		//validate the token
		token, err := validateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission Denied!"))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission Denied!"))
			return
		}
		
		handler.ServeHTTP(w, r)
	})
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}
