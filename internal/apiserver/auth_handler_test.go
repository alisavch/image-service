package apiserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/service/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_signUp(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Test with correct values",
			inputBody: `{"username": "username", "password": "12345"}`,
			inputUser: models.User{
				Username: "username",
				Password: "12345",
			},
			fn: func(mockAuthorization *mocks.Authorization, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := new(mocks.Authorization)
			tt.fn(auth, tt.inputUser)

			services := &service.Service{Authorization: auth}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/sign-up", s.signUp()).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/sign-up",
				bytes.NewBufferString(tt.inputBody))
			s.ServeHTTP(w, req)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_singIn(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.Authorization, username string, password string)
	type user struct {
		username string
		password string
	}

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            user
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Test with correct values",
			inputBody: `{"username": "username", "password": "12345"}`,
			inputUser: user{username: "username", password: "12345"},
			fn: func(mockAuthorization *mocks.Authorization, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "\"token\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := new(mocks.Authorization)
			tt.fn(auth, tt.inputUser.username, tt.inputUser.password)

			services := &service.Service{Authorization: auth}
			s := Server{router: mux.NewRouter(), service: services}

			s.router.HandleFunc("/api/sign-in", s.signIn()).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/sign-in",
				bytes.NewBufferString(tt.inputBody))
			s.ServeHTTP(w, req)
			require.Equal(t, tt.expectedStatusCode, w.Code)
			require.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
