package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/model"
	"github.com/ynuraddi/t-medods/repository"
	"golang.org/x/crypto/bcrypt"
)

type payload struct {
	ID        string    `json:"jti"`
	Subject   string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func (p *payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return model.ErrExpiredToken
	}
	return nil
}

type authService struct {
	accessKey  []byte
	refreshKey []byte
	repo       repository.ISessionRepository
}

func NewAuthService(config *config.Config, repo repository.ISessionRepository) *authService {
	return &authService{
		accessKey:  []byte(config.TokenAccessKey),
		refreshKey: []byte(config.TokenRefreshKey),
		repo:       repo,
	}
}

func (s *authService) CreateSession(ctx context.Context, userID string) (access, refresh string, err error) {
	tokenID, err := s.createTokenID()
	if err != nil {
		return "", "", err
	}

	access, err = s.createToken(userID, tokenID, 5*time.Minute, s.accessKey)
	if err != nil {
		return "", "", fmt.Errorf("failed generate access token: %w", err)
	}

	refresh, err = s.createToken(userID, tokenID, 1*time.Hour, s.refreshKey)
	if err != nil {
		return "", "", fmt.Errorf("failed generate refresh token: %w", err)
	}

	hashRefresh, err := s.hashToken(refresh)
	if err != nil {
		return "", "", fmt.Errorf("failed hash refresh token: %w", err)
	}

	if err = s.repo.CreateSession(ctx, model.Session{
		UserID:    userID,
		TokenHash: hashRefresh,
	}); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) RefreshSession(ctx context.Context, refreshToken string) (access, refresh string, err error) {
	payload, err := s.verifyToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	session, err := s.repo.SessionByUser(ctx, payload.Subject)
	if err != nil {
		return "", "", err
	}

	if !s.compareTokensHash(session.TokenHash, refreshToken) {
		return "", "", model.ErrExpiredToken
	}

	access, refresh, err = s.CreateSession(ctx, payload.Subject)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *authService) createToken(userID string, tokenID string, duration time.Duration, key []byte) (string, error) {
	iat := time.Now()
	exp := iat.Add(duration)

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": iat.Unix(),
		"exp": exp.Unix(),
		"jti": tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(key)
}

func (s *authService) verifyToken(token string) (*payload, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &payload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrInvalidToken
		}
		return s.refreshKey, nil
	})
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, model.ErrExpiredToken) {
			return nil, model.ErrExpiredToken
		}
		return nil, model.ErrInvalidToken
	}

	payload, ok := parsedToken.Claims.(*payload)
	if !ok {
		return nil, model.ErrInvalidToken
	}

	return payload, nil
}

func (s *authService) createRefresh() (string, error) {
	jwtID := make([]byte, 72)
	_, err := rand.Read(jwtID)
	if err != nil {
		return "", fmt.Errorf("failed generate refresh token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(jwtID), nil
}

func (s *authService) hashToken(token string) (string, error) {
	hashedtoken, err := bcrypt.GenerateFromPassword([]byte(token), 32)
	if err != nil {
		return "", err
	}
	return string(hashedtoken), nil
}

func (s *authService) compareTokensHash(hashedtoken, plaintoken string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedtoken), []byte(plaintoken)) == nil
}
