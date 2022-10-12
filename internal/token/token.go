package token

import "time"


// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken generates a new token for a specific username valid for a specific duration of time
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks and verifies the payload of the token
	VerifyToken(token string) (*Payload, error)
}