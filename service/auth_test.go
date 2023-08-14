package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/model"
	mock_repository "github.com/ynuraddi/t-medods/repository/mock"
)

var errUnexpected error = errors.New("unexpected")

func TestCreateSession(t *testing.T) {
	testCases := []struct {
		name        string
		buildStubs  func(repo *mock_repository.MockISessionRepository)
		checkReturn func(access, refresh string, err error)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mock_repository.MockISessionRepository) {
				repo.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkReturn: func(access, refresh string, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, access)
				require.NotEmpty(t, refresh)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(repo *mock_repository.MockISessionRepository) {
				repo.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(errUnexpected)
			},
			checkReturn: func(access, refresh string, err error) {
				require.Error(t, err)
				require.Empty(t, access)
				require.Empty(t, refresh)
			},
		},
	}

	for _, test := range testCases {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		mockRepo := mock_repository.NewMockISessionRepository(ctrl)

		service := authService{
			secretKey: []byte("123"),
			repo:      mockRepo,
		}

		test.buildStubs(mockRepo)
		test.checkReturn(service.CreateSession(context.Background(), "123"))
	}
}

func TestCreateToken(t *testing.T) {
	service := NewAuthService(&config.Config{
		TokenAccessKey: "123",
	}, nil)

	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := service.createToken("1", duration, service.secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := service.verifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Subject, "1")
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpireToken(t *testing.T) {
	service := NewAuthService(&config.Config{
		TokenAccessKey: "123",
	}, nil)

	duration := time.Minute

	token, err := service.createToken("1", -duration, service.secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := service.verifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, model.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidAlgToken(t *testing.T) {
	service := NewAuthService(&config.Config{
		TokenAccessKey: "123",
	}, nil)

	payload := &payload{
		Subject:   "1",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Minute),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = service.verifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, model.ErrInvalidToken.Error())
	require.Nil(t, payload)
}
