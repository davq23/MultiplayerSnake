package middlewares

import (
	"context"
	"davidmultiplayersnake/api/config"
	"davidmultiplayersnake/api/security"
	"errors"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JWTMiddleware wraps a handlerfunc and blocks it to allow JWT authentication
func JWTMiddleware(next http.HandlerFunc, authRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("papoman")

		if authRequired && (err != nil || cookie == nil) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if cookie != nil && authRequired {
			token, err := security.VerifyToken(cookie.Value, config.JWTSecret)

			if err != nil || !token.Valid {
				log.Println(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			playerClaims, err := decipherClaims(token)

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

func decipherClaims(token *jwt.Token) (security.PlayerClaims, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return security.PlayerClaims{}, errors.New("Unable to parse claims")
	}

	playerClaims := security.PlayerClaims{
		HubName:    mapClaims["hub_name"].(string),
		PlayerName: mapClaims["player_name"].(string),
	}

	playerClaims.StandardClaims.ExpiresAt = int64(mapClaims["exp"].(float64))
	return playerClaims, nil
}
