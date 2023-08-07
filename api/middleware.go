package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"log"
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

func getUserDetailsFromClaims(r *http.Request) map[string]string {
	tokenString := r.Header.Get("access-token")
	token, _ := validateToken(tokenString)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("error when retrieving claims")
		return nil
	}

	userDetailsMap := make(map[string]string)
	usr := claims["UserDetails"].(map[string]interface{})
	for key, value := range usr {
		userDetailsMap[key] = value.(string)
	}

	return userDetailsMap
}
