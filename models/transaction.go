package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        primitive.ObjectID    `bson:"_id,omitempty"`
	UserID    primitive.ObjectID    `bson:"user_id"`
	Type      string    			`bson:"type"` // withdrawal or deposit
	Amount    float64   			`bson:"amount"`
	CreatedAt time.Time 			`bson:"cretaed_at"`
	UpdatedAt time.Time 			`bson:"updated_at"`
}