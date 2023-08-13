package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-medods/model"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) refresh(c echo.Context) error {
	var req refreshRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind refresh request", err)
		c.JSON(http.StatusUnprocessableEntity, errorResponce(err))
	}
	if req.RefreshToken == "" {
		err := fmt.Errorf("bad request")
		s.logger.Error("refresh", err)
		c.JSON(http.StatusUnprocessableEntity, errorResponce(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	access, refresh, err := s.service.Auth.RefreshSession(ctx, req.RefreshToken)
	if errors.Is(err, model.ErrInvalidToken) {
		s.logger.Error("invalid refresh token", err)
		return c.JSON(http.StatusUnauthorized, err)
	}

	if errors.Is(err, model.ErrExpiredToken) {
		s.logger.Error("expired refresh token", err)
		return c.JSON(http.StatusUnauthorized, err)
	}

	if err != nil {
		s.logger.Error("failed refresh tokens", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, loginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
