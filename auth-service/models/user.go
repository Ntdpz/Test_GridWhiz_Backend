package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the database
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // "-" means don't return in JSON
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// EndToken represents invalid/expired tokens
type EndToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Token     string             `bson:"token" json:"token"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
