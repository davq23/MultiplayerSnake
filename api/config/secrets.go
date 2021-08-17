package config

// JWTSecret is the secret used to verify the JWT token
const JWTSecret = os.Getenv("jwtSecret")
