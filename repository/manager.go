package repository

import (
	"context"

	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/model"
	"github.com/ynuraddi/t-medods/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	UserByName(ctx context.Context, username string) (dbuser model.User, err error)
}

type Manager struct {
	User IUserRepository
}

func New(config config.Config, logger logger.Logger, client *mongo.Client) *Manager {
	dbmedods := client.Database(config.MongoDBName)

	userRepostiory := mongodb.NewUserRepostiory(dbmedods)

	repository := &Manager{
		User: userRepostiory,
	}
	logger.Info("repository success inited")
	return repository
}
