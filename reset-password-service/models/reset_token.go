package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResetToken represents a password reset token
type ResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
	IsUsed    bool               `bson:"is_used"`
	CreatedAt time.Time          `bson:"created_at"`
}
