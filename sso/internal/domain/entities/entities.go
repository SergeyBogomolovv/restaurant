package entities

import "time"

type RefreshToken struct {
	EntityID  string    `json:"entity_id"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}
