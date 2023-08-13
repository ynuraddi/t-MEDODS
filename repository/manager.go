package repository

import (
	"context"

	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/model"
	"github.com/ynuraddi/t-medods/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ISessionRepository interface {
	CreateSession(ctx context.Context, sess model.Session) error
	SessionByUser(ctx context.Context, userID string) (session model.Session, err error)
}

type Manager struct {
	Sess ISessionRepository
}

func New(config config.Config, logger logger.Logger, client *mongo.Client) *Manager {
	dbmedods := client.Database(config.MongoDBName)

	sessionRepostiory := mongodb.NewSessionRepostiory(dbmedods)

	repository := &Manager{
		Sess: sessionRepostiory,
	}
	logger.Info("repository success inited")
	return repository
}
