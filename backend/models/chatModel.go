package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chat represents a chat in MongoDB
type Chat struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	ChatName      string               `bson:"chatName,omitempty"`
	IsGroupChat   bool                 `bson:"isGroupChat"`
	Users         []primitive.ObjectID `bson:"users,omitempty"`
	LatestMessage primitive.ObjectID   `bson:"latestMessage,omitempty"`
}
