package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/model"
	"github.com/ynuraddi/t-medods/service"
	mock_service "github.com/ynuraddi/t-medods/service/mock"
)

func TestRefresh(t *testing.T) {
	testCases := []struct {
		name          string
		authToken     string
		refrToken     string
		buildStubs    func(service *mock_service.MockIAuthService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			authToken: "Bearer 123",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("123", "123", nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NoAuthHeader",
			authToken: "",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "InvalidAuthHeader",
			authToken: "Bearer ",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoRefreshToken",
			authToken: "Bearer 123",
			refrToken: `{"refresh_token":""}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidJSON",
			authToken: "Bearer 123",
			refrToken: `{"refres\\}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:      "TokenInvalid",
			authToken: "Bearer 123",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("", "", model.ErrInvalidToken)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

				type errorResponse struct {
					Error string `josn:"error"`
				}
				var get errorResponse
				err := json.NewDecoder(recorder.Body).Decode(&get)
				require.NoError(t, err)
				require.Equal(t, get.Error, model.ErrInvalidToken.Error())
			},
		},
		{
			name:      "TokenExpired",
			authToken: "Bearer 123",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("", "", model.ErrExpiredToken)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

				type errorResponse struct {
					Error string `josn:"error"`
				}
				var get errorResponse
				err := json.NewDecoder(recorder.Body).Decode(&get)
				require.NoError(t, err)
				require.Equal(t, get.Error, model.ErrExpiredToken.Error())
			},
		},
		{
			name:      "InternalError",
			authToken: "Bearer 123",
			refrToken: `{"refresh_token":"1234567890"}`,
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().RefreshSession(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return("", "", errUnexpected)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessService := mock_service.NewMockIAuthService(ctrl)
			test.buildStubs(sessService)

			service := &service.Manager{Auth: sessService}
			server := Server{
				logger:  logger.NewLogger(&config.Config{LogLevel: 0}, nil),
				service: service,
				router:  echo.New(),
			}

			req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(test.refrToken))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, test.authToken)

			rec := httptest.NewRecorder()

			c := server.router.NewContext(req, rec)
			server.refresh(c)

			test.checkResponce(t, rec)
		})
	}
}
