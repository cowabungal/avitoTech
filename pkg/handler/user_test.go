package handler

import (
	"avitoTech"
	"avitoTech/pkg/service"
	service_mocks "avitoTech/pkg/service/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_Balance(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockUser, user avitoTech.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            avitoTech.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id": 1}`,
			inputUser: avitoTech.User{
				UserId: 1,
			},
			mockBehavior: func(r *service_mocks.MockUser, user avitoTech.User) {
				ans := avitoTech.User{
					UserId:  1,
					Balance: 100,
				}
				r.EXPECT().Balance(user.UserId).Return(&ans, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":100}`,
		},
		{
			name:      "Bad",
			inputBody: `{"bad": 1}`,
			inputUser: avitoTech.User{},
			mockBehavior: func(r *service_mocks.MockUser, user avitoTech.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "no user",
			inputBody: `{"user_id":2}`,
			inputUser: avitoTech.User{
				UserId: 2,
			},
			mockBehavior: func(r *service_mocks.MockUser, user avitoTech.User) {
				err := errors.New("no user")
				r.EXPECT().Balance(user.UserId).Return(nil, err)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user has no balance"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.GET("/balance", handler.Balance)

			// Create Request
			w := httptest.NewRecorder()
			req:= httptest.NewRequest("GET", "/balance",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_TopUp(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockUser, input Input)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            Input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id":1,"amount":10}`,
			inputUser: Input{
				UserId: 1,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {
				ans := avitoTech.User{
					UserId:  1,
					Balance: 10,
				}
				r.EXPECT().TopUp(input.UserId, input.Amount).Return(&ans, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":10}`,
		},
		{
			name:      "Bad",
			inputBody: `{"user_id":1}`,
			inputUser: Input{
				UserId: 1,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "internal error",
			inputBody: `{"user_id":1,"amount":10}`,
			inputUser: Input{
				UserId: 1,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {
				err := errors.New("internal error")
				r.EXPECT().TopUp(input.UserId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/top-up", handler.TopUp)

			// Create Request
			w := httptest.NewRecorder()
			req:= httptest.NewRequest("POST", "/top-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_Debit(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockUser, input Input)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            Input
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id":1,"amount":10}`,
			inputUser: Input{
				UserId: 1,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {
				ans := avitoTech.User{
					UserId:  1,
					Balance: 0,
				}
				r.EXPECT().Debit(input.UserId, input.Amount).Return(&ans, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":0}`,
		},
		{
			name:      "Bad",
			inputBody: `{"user_id":1}`,
			inputUser: Input{
				UserId: 1,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "no balance",
			inputBody: `{"user_id":1,"amount":10}`,
			inputUser: Input{
				UserId: 1,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {
				err := errors.New("insufficient funds")
				r.EXPECT().Debit(input.UserId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"insufficient funds"}`,
		},
		{
			name:      "internal error",
			inputBody: `{"user_id":1,"amount":10}`,
			inputUser: Input{
				UserId: 1,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Input) {
				err := errors.New("internal error")
				r.EXPECT().Debit(input.UserId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/debit", handler.Debit)

			// Create Request
			w := httptest.NewRecorder()
			req:= httptest.NewRequest("POST", "/debit",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_Transfer(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockUser, input Transfer)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            Transfer
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id":1,"to_id":2,"amount":10}`,
			inputUser: Transfer{
				UserId: 1,
				ToId: 2,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Transfer) {
				ans := avitoTech.User{
					UserId:  2,
					Balance: 10,
				}
				r.EXPECT().Transfer(input.UserId, input.ToId, input.Amount).Return(&ans, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":2,"balance":10}`,
		},
		{
			name:      "Bad",
			inputBody: `{"user_id":1}`,
			inputUser: Transfer{
				UserId: 1,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Transfer) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "no balance",
			inputBody: `{"user_id":1,"to_id":2,"amount":10}`,
			inputUser: Transfer{
				UserId: 1,
				ToId: 2,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Transfer) {
				err := errors.New("insufficient funds")
				r.EXPECT().Transfer(input.UserId, input.ToId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"insufficient funds"}`,
		},
		{
			name:      "no user",
			inputBody: `{"user_id":1,"to_id":2,"amount":10}`,
			inputUser: Transfer{
				UserId: 1,
				ToId: 2,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Transfer) {
				err := errors.New("the recipient has no balance")
				r.EXPECT().Transfer(input.UserId, input.ToId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"the recipient has no balance"}`,
		},
		{
			name:      "internal error",
			inputBody: `{"user_id":1, "to_id":2, "amount":10}`,
			inputUser: Transfer{
				UserId: 1,
				ToId: 2,
				Amount: 10,
			},
			mockBehavior: func(r *service_mocks.MockUser, input Transfer) {
				err := errors.New("internal error")
				r.EXPECT().Transfer(input.UserId, input.ToId, input.Amount).Return(nil, err)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/transfer", handler.Transfer)

			// Create Request
			w := httptest.NewRecorder()
			req:= httptest.NewRequest("POST", "/transfer",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
