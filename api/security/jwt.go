package security

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// GetToken creates a new access token that expires in timeout minutes
func GetToken(playerName, hubname, secret string, playerScore int, timeout int64) (token string, err error) {
	// Auth claims
	claims := PlayerClaims{}
	claims.HubName = hubname
	claims.PlayerName = playerName
	claims.PlayerScore = playerScore
	claims.StandardClaims.ExpiresAt = timeout

	// Sign token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = at.SignedString([]byte(secret))

	return
}

// VerifyToken verifies the JWT token
func VerifyToken(tokenString string, tokenSecret string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	return
}
