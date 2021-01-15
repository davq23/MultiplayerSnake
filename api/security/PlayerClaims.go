package security

import (
	"github.com/dgrijalva/jwt-go"
)

// PlayerClaimsType is an identifier to save player claims into request context
const PlayerClaimsType = 0

// PlayerClaims are the claims to identify a player in a hub
type PlayerClaims struct {
	HubName    string `json:"hub_name"`
	PlayerName string `json:"player_name"`
	jwt.StandardClaims
}
