package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a message in MongoDB
type Message struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Sender    primitive.ObjectID   `bson:"sender,omitempty"`
	Content   string               `bson:"content,omitempty"`
	Chat      primitive.ObjectID   `bson:"chat,omitempty"`
	ReadBy    []primitive.ObjectID `bson:"readBy,omitempty"`
	CreatedAt time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty"`
}
