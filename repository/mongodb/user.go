package mongodb

import (
	"context"
	"errors"

	"github.com/ynuraddi/t-medods/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepostiory struct {
	coll *mongo.Collection
}

func NewUserRepostiory(client *mongo.Database) *userRepostiory {
	return &userRepostiory{
		coll: client.Collection("users"),
	}
}

func (r *userRepostiory) UserByName(ctx context.Context, username string) (dbuser model.User, err error) {
	if err = r.coll.FindOne(ctx, bson.M{"username": username}).Decode(dbuser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrNoUser
		}

		return model.User{}, err
	}
	return
}
