package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var invalidFormatError string = "Invalid token format provided."
var invalidTokenError string = "Invalid token provided."

func Mount(next http.Handler) http.Handler {
	token := os.Getenv("AUTHORIZATION_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			providedToken, err := extractTokenFromRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}

			if providedToken != token {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(invalidTokenError))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func extractTokenFromRequest(r *http.Request) (string, error) {
	providedToken := r.Header.Get("Authorization")
	s := strings.Split(providedToken, " ")
	if len(s) != 2 {
		return "", fmt.Errorf(invalidFormatError)
	}

	return s[1], nil
}
