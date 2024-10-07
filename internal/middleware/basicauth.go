package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(auth, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
			return
		}

		if subtle.ConstantTimeCompare([]byte(pair[0]), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pair[1]), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
