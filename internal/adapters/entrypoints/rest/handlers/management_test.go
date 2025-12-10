package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// nolint:funlen
func TestManagementInterfaceHandlerHappyPath(t *testing.T) {
	t.Run("should render login template when not logged in with credentials not set", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: false,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		// Return error to simulate no session (not logged in)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
			Name: liquidity_provider.ManagementLoginTemplate,
			Data: liquidity_provider.ManagementTemplateData{
				CredentialsSet: false,
				BaseUrl:        "http://localhost:8080",
			},
		}, nil)

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		body := recorder.Body.String()

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "text/html", recorder.Header().Get("Content-Type"))

		// Verify login template structure
		assert.Contains(t, body, "<title>Management Login</title>")
		assert.Contains(t, body, `id="login-form"`)
		assert.Contains(t, body, `data-testid="login-username-input"`)
		assert.Contains(t, body, `data-testid="login-password-input"`)

		// When CredentialsSet is false, new credential fields should appear
		assert.Contains(t, body, `id="new-username"`)
		assert.Contains(t, body, `id="new-password"`)

		mockStore.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should render login template when not logged in with credentials already set", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: false,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
			Name: liquidity_provider.ManagementLoginTemplate,
			Data: liquidity_provider.ManagementTemplateData{
				CredentialsSet: true, // Credentials already set
				BaseUrl:        "http://localhost:8080",
			},
		}, nil)

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		body := recorder.Body.String()

		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verify login template structure
		assert.Contains(t, body, "<title>Management Login</title>")
		assert.Contains(t, body, `id="login-form"`)
		assert.Contains(t, body, `data-testid="login-username-input"`)
		assert.Contains(t, body, `data-testid="login-password-input"`)

		// When CredentialsSet is true, new credential fields should NOT appear
		assert.NotContains(t, body, `id="new-username"`)
		assert.NotContains(t, body, `id="new-password"`)

		mockStore.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})
}

func TestManagementInterfaceHandlerErrorCases(t *testing.T) {
	t.Run("should render error template when use case fails", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: false,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(nil, errors.New("database error"))

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// Error template should be rendered (still returns 200 but with error content)
		assert.Equal(t, http.StatusOK, recorder.Code)

		mockStore.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})
}

// nolint:funlen
func TestManagementInterfaceHandlerSecurityHeaders(t *testing.T) {
	t.Run("should set security headers when enabled", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: true,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
			Name: liquidity_provider.ManagementLoginTemplate,
			Data: liquidity_provider.ManagementTemplateData{
				CredentialsSet: false,
				BaseUrl:        "http://localhost:8080",
			},
		}, nil)

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		// Verify security headers are set
		assert.NotEmpty(t, recorder.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "max-age=63072000; includeSubDomains; preload", recorder.Header().Get("Strict-Transport-Security"))
		assert.Equal(t, "DENY", recorder.Header().Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", recorder.Header().Get("X-Content-Type-Options"))

		mockStore.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should not set security headers when disabled", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: false,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
			Name: liquidity_provider.ManagementLoginTemplate,
			Data: liquidity_provider.ManagementTemplateData{
				CredentialsSet: false,
				BaseUrl:        "http://localhost:8080",
			},
		}, nil)

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		// Verify security headers are NOT set
		assert.Empty(t, recorder.Header().Get("Content-Security-Policy"))
		assert.Empty(t, recorder.Header().Get("Strict-Transport-Security"))
		assert.Empty(t, recorder.Header().Get("X-Frame-Options"))
		assert.Empty(t, recorder.Header().Get("X-Content-Type-Options"))

		mockStore.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})
}

func TestManagementInterfaceHandlerResponseFormat(t *testing.T) {
	t.Run("should return text/html content type", func(t *testing.T) {
		env := environment.ManagementEnv{
			EnableSecurityHeaders: false,
		}

		request := httptest.NewRequest(http.MethodGet, "/management", nil)
		recorder := httptest.NewRecorder()

		mockStore := new(mocks.StoreMock)
		mockStore.On("Get", mock.Anything, "lp-session").Return(nil, errors.New("no session"))

		mockUseCase := new(mocks.GetManagementUiDataUseCaseMock)
		mockUseCase.On("Run", mock.Anything, false).Return(&liquidity_provider.ManagementTemplate{
			Name: liquidity_provider.ManagementLoginTemplate,
			Data: liquidity_provider.ManagementTemplateData{
				CredentialsSet: false,
				BaseUrl:        "http://localhost:8080",
			},
		}, nil)

		handlerFunc := handlers.NewManagementInterfaceHandler(env, mockStore, mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "text/html", recorder.Header().Get("Content-Type"))
	})
}
