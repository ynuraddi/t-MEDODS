package mongodb

import (
	"context"
	"fmt"

	"github.com/ynuraddi/t-medods/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sessionRepostiory struct {
	coll *mongo.Collection
}

func NewSessionRepostiory(client *mongo.Database) *sessionRepostiory {
	return &sessionRepostiory{
		coll: client.Collection("sessions"),
	}
}

func (r *sessionRepostiory) CreateSession(ctx context.Context, sess model.Session) error {
	filter := bson.M{"user_id": sess.UserID}
	update := bson.M{"$set": bson.M{
		"user_id":    sess.UserID,
		"token_hash": sess.TokenHash,
	}}
	opts := options.Update().SetUpsert(true)

	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed upsert session: %w", err)
	}

	return nil
}

func (r *sessionRepostiory) SessionByUser(ctx context.Context, userID string) (session model.Session, err error) {
	filter := bson.M{"user_id": userID}
	if err = r.coll.FindOne(ctx, filter).Decode(&session); err != nil {
		return model.Session{}, fmt.Errorf("failed find session: %w", err)
	}

	return session, nil
}
