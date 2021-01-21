package security

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// PlayerClaimsType is an identifier to save player claims into request context
const PlayerClaimsType = 0

// PlayerClaims are the claims to identify a player in a hub
type PlayerClaims struct {
	HubName     string `json:"hub_name"`
	PlayerName  string `json:"player_name"`
	PlayerScore int    `json:"player_score"`
	jwt.StandardClaims
}

func DecipherClaims(token *jwt.Token) (PlayerClaims, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return PlayerClaims{}, errors.New("Unable to parse claims")
	}

	playerClaims := PlayerClaims{
		HubName:     mapClaims["hub_name"].(string),
		PlayerName:  mapClaims["player_name"].(string),
		PlayerScore: mapClaims["player_name"].(int),
	}

	playerClaims.StandardClaims.ExpiresAt = int64(mapClaims["exp"].(float64))
	return playerClaims, nil
}
