package service

import (
	"context"

	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/repository"
)

type IAuthService interface {
	CreateSession(ctx context.Context, uid string) (access, refresh string, err error)
	RefreshSession(ctx context.Context, refreshToken string) (access, refresh string, err error)
}

type Manager struct {
	Auth IAuthService
}

func New(config *config.Config, logger logger.Logger, repo *repository.Manager) *Manager {
	authService := NewAuthService(config, repo.Sess)

	service := &Manager{
		Auth: authService,
	}
	return service
}
