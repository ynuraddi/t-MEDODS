package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/service"
)

type Server struct {
	config *config.Config
	logger logger.Logger

	service *service.Manager
	router  *echo.Echo
}

const shutdownTimeout = 5 * time.Second

func NewServer(config *config.Config, logger logger.Logger, service *service.Manager) *Server {
	server := Server{
		config:  config,
		logger:  logger,
		service: service,
	}

	router := echo.New()
	router.Server.ReadTimeout = 10 * time.Second
	router.Server.IdleTimeout = 30 * time.Second
	server.router = router
	server.setupRouter()

	return &server
}

func (s *Server) setupRouter() {
	router := s.router

	router.POST("/auth/:id", s.login)
	router.POST("/refresh", s.refresh)
}

func (s *Server) Serve(ctx context.Context) error {
	msg := fmt.Sprintf("starting serve on %s:%d good job!", s.config.HttpHost, s.config.HttpPort)
	s.logger.Info(msg)

	addr := fmt.Sprintf("%s:%d", s.config.HttpHost, s.config.HttpPort)
	go func() {
		if err := s.router.Start(addr); err != http.ErrServerClosed {
			s.logger.Error("server stoped", err)
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer func() {
		cancel()
	}()

	s.logger.Info("gracefull shutdown server is started")
	if err := s.router.Shutdown(ctxShutdown); err != nil {
		log.Fatalln("gracefull shutdown is failed", err)
	}
	s.logger.Info("gracefull shutdown server is complete")

	return nil
}
