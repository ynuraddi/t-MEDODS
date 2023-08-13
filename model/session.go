package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	TokenHash string             `bson:"token_hash"`
}
