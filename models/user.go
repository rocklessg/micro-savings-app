package models

import "time"

// User represents a registered user
type User struct {
	ID              string    `bson:"_id,omitempty"`
	Name            string    `bson:"name"`
	Email           string    `bson:"email"`
	PasswordHash    string    `bson:"password_hash"`
	SavingsBalance  float64   `bson:"savings_balance"`
	InvestmentBalance float64 `bson:"investment_balance"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}
