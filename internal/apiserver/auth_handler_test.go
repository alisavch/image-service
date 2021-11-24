package apiserver

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alisavch/image-service/internal/apiserver/mocks"
	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_signUp(t *testing.T) {
	type fnBehavior func(mockAuthorization *mocks.ServiceOperations, user models.User)

	asString := "00000000-0000-0000-0000-000000000000"
	s := uuid.MustParse(asString)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		fn                   fnBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Create user without errors",
			inputBody: `{"username": "username", "password": "12345"}`,
			inputUser: models.User{
				Username: "username",
				Password: "12345",
			},
			fn: func(mockAuthorization *mocks.ServiceOperations, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, user).Return(s, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "\"User registered successfully\"\n",
		},
		{
			name:      "Password must not be empty",
			inputBody: `{"username": "uuu", "password": ""}`,
			inputUser: models.User{
				Username: "username",
				Password: "12345",
			},
			fn: func(mockAuthorization *mocks.ServiceOperations, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, user).Return(s, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"error\":\"password must not be empty\"}\n",
		},
		{
			name:      "Auth header is empty",
			inputBody: `{"username": "uuu", "password": "ppp"}`,
			inputUser: models.User{
				Username: "uuu",
				Password: "ppp",
			},
			fn: func(mockAuthorization *mocks.ServiceOperations, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, mock.Anything).Return(s, utils.ErrEmptyHeader)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"error\":\"auth header is empty\"}\n",
		},
		{
			name:      "Username must not be empty",
			inputBody: `{"username": "", "password": "ppp"}`,
			inputUser: models.User{
				Username: "uuu",
				Password: "ppp",
			},
			fn: func(mockAuthorization *mocks.ServiceOperations, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, mock.Anything).Return(s, utils.ErrEmptyUsername)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"error\":\"username must not be empty\"}\n",
		},
		{
			name:      "EOF",
			inputBody: "",
			inputUser: models.User{
				Username: "uuu",
				Password: "ppp",
			},
			fn: func(mockAuthorization *mocks.ServiceOperations, user models.User) {
				mockAuthorization.On("CreateUser", mock.Anything, mock.Anything).Return(s, fmt.Errorf("EOF"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"error\":\"EOF\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			mq := broker.NewAMQPBroker(mockSO, mockBucket)

			s := NewServer(mq, currentService)

			tt.fn(mockSO, tt.inputUser)

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
	type fnBehavior func(mockAuthorization *mocks.ServiceOperations, username string, password string)
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
			name:      "Generate token without errors",
			inputBody: `{"username": "username", "password": "12345"}`,
			inputUser: user{username: "username", password: "12345"},
			fn: func(mockAuthorization *mocks.ServiceOperations, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "\"token\"\n",
		},
		{
			name:      "Password must not be empty",
			inputBody: `{"username": "uuu", "password": ""}`,
			fn: func(mockAuthorization *mocks.ServiceOperations, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("", utils.ErrEmptyPassword)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"password must not be empty\"}\n",
		},
		{
			name:      "Auth header is empty",
			inputBody: `{"username": "uuu", "password": "ppp"}`,
			inputUser: user{username: "uuu", password: "ppp"},
			fn: func(mockAuthorization *mocks.ServiceOperations, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("", utils.ErrEmptyHeader)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"auth header is empty\"}\n",
		},
		{
			name:      "Username must not be empty",
			inputBody: `{"username": "", "password": "12345"}`,
			inputUser: user{username: "username", password: "12345"},
			fn: func(mockAuthorization *mocks.ServiceOperations, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("token", nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"username must not be empty\"}\n",
		},
		{
			name:      "EOF",
			inputBody: "",
			inputUser: user{username: "username", password: "12345"},
			fn: func(mockAuthorization *mocks.ServiceOperations, username string, password string) {
				mockAuthorization.On("GenerateToken", mock.Anything, username, password).Return("token", nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"error\":\"EOF\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBucket := new(mocks.S3Bucket)
			mockSO := new(mocks.ServiceOperations)

			currentService := NewAPI(mockSO, mockBucket)
			mq := broker.NewAMQPBroker(mockSO, mockBucket)

			s := NewServer(mq, currentService)

			tt.fn(mockSO, tt.inputUser.username, tt.inputUser.password)

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
