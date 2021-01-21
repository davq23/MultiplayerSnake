package middlewares

import (
	"context"
	"davidmultiplayersnake/api/config"
	"davidmultiplayersnake/api/security"
	"log"
	"net/http"
	"strings"
)

// JWTMiddleware wraps a handlerfunc and blocks it to allow JWT authentication
func JWTMiddleware(next http.HandlerFunc, authRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearToken := r.Header.Get("Authorization")

		strArr := strings.Split(bearToken, " ")

		if authRequired && (len(strArr) != 2 || strArr[0] != "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if authRequired {
			token, err := security.VerifyToken(strArr[1], config.JWTSecret)

			if err != nil || !token.Valid {
				log.Println(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			playerClaims, err := security.DecipherClaims(token)

			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			next(w, r.WithContext(context.WithValue(r.Context(), security.PlayerClaimsType, playerClaims)))
			return
		}

		next(w, r)
	}
}

// WSJWTMiddleware wraps a handlerfunc and blocks it to allow JWT authentication in WebSockets
func WSJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["user_info"]

		if !ok || len(values) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := security.VerifyToken(values[0], config.JWTSecret)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		playerClaims, err := security.DecipherClaims(token)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), security.PlayerClaimsType, playerClaims)))
	}
}
