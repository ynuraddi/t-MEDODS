package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-medods/model"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) refresh(c echo.Context) error {
	access := c.Request().Header.Get(echo.HeaderAuthorization)
	if len(access) == 0 {
		err := fmt.Errorf("empty access token in Authorization Header")
		s.logger.Error("", err)
		return c.JSON(http.StatusUnauthorized, errorResponce(err))
	}

	access = strings.TrimPrefix(access, "Bearer ")
	if len(access) == 0 {
		err := fmt.Errorf("invalid access token in Authorization Header")
		s.logger.Error("", err)
		return c.JSON(http.StatusUnauthorized, errorResponce(err))
	}

	var req refreshRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind refresh request", err)
		return c.JSON(http.StatusUnprocessableEntity, errorResponce(err))
	}
	if len(req.RefreshToken) == 0 {
		err := fmt.Errorf("bad request")
		s.logger.Error("refresh", err)
		return c.JSON(http.StatusBadRequest, errorResponce(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	access, refresh, err := s.service.Auth.RefreshSession(ctx, access, req.RefreshToken)
	if errors.Is(err, model.ErrInvalidToken) {
		s.logger.Error("", err)
		return c.JSON(http.StatusUnauthorized, errorResponce(err))
	}
	if errors.Is(err, model.ErrExpiredToken) {
		s.logger.Error("", err)
		return c.JSON(http.StatusUnauthorized, errorResponce(err))
	}
	if err != nil {
		s.logger.Error("failed refresh tokens", err)
		return c.JSON(http.StatusInternalServerError, errorResponce(err))
	}

	return c.JSON(http.StatusOK, refreshResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
