package transport

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/service"
	mock_service "github.com/ynuraddi/t-medods/service/mock"
)

var errUnexpected error = errors.New("unexpected")

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(service *mock_service.MockIAuthService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return("123", "123", nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "InvalidParam",
			param: "",
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "InternalError",
			param: "1",
			buildStubs: func(service *mock_service.MockIAuthService) {
				service.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return("", "", errUnexpected)
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

			req := httptest.NewRequest(http.MethodPost, "/auth/", nil)
			rec := httptest.NewRecorder()

			c := server.router.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.param)
			server.login(c)

			test.checkResponce(t, rec)
		})
	}
}
