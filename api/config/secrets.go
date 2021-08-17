package config

import "os"

// JWTSecret is the secret used to verify the JWT token
var JWTSecret = os.Getenv("jwtSecret")
