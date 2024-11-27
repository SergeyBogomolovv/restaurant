package entities

import "time"

type CustomerEntity struct {
	CustomerID   string    `db:"customer_id"`
	Name         string    `db:"name"`
	BirthDate    time.Time `db:"birth_date"`
	TotalSpent   float64   `db:"total_spent"`
	Email        string    `db:"email"`
	Password     []byte    `db:"password"`
	RegisteredAt time.Time `db:"registered_at"`
}

type WaiterEntity struct {
	WaiterID    string     `db:"waiter_id"`
	Login       string     `db:"login"`
	Password    []byte     `db:"password"`
	FirstName   string     `db:"first_name"`
	LastName    string     `db:"last_name"`
	HiredAt     time.Time  `db:"hired_at"`
	FiredAt     *time.Time `db:"fired_at"`
	FiredReason *string    `db:"fired_reason"`
	Rating      float64    `db:"rating"`
}

type AdminEntity struct {
	AdminID  string  `db:"admin_id"`
	Note     *string `db:"note"`
	Login    string  `db:"login"`
	Password []byte  `db:"password"`
}

type RefreshTokenEntity struct {
	EntityID  string    `json:"entity_id"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}
