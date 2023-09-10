package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in MongoDB
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

// User login information
type UserLoginInformation struct {
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
}
