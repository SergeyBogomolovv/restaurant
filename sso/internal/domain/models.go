package domain

import (
	"time"
)

type RefreshTokenPayload struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Role      string    `json:"role"`
}
