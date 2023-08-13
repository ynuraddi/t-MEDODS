package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) login(c echo.Context) error {
	userID := c.Param("id")

	if userID == "" {
		err := fmt.Errorf("auth/:id, id param is required")
		s.logger.Error("invalid param", err)
		return c.JSON(http.StatusBadRequest, errorResponce(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	access, refresh, err := s.service.Auth.CreateSession(ctx, userID)
	if err != nil {
		s.logger.Error("failed login user", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, loginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
