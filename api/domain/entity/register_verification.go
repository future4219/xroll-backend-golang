package entity

import "time"

type RegisterVerification struct {
	RegisterVerificationID   string
	Email                    string
	HashedPassword           string
	HashedAuthenticationCode string
	ExpiresAt                time.Time
}
