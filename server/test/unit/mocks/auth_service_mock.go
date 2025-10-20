package mocks

import (
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

var _ authfx.AuthServiceInterface = (*MockAuthService)(nil)

func (m *MockAuthService) GetRegistrationMail(email string) error {
	args := m.Called(email)

	return args.Error(0)
}

func (m *MockAuthService) Register(registrationTokenString string, body *authfx.RegisterBody) (*models.PublicUser, string, error) {
	args := m.Called(registrationTokenString, body)

	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}

	return args.Get(0).(*models.PublicUser), args.String(1), args.Error(2)
}

func (m *MockAuthService) Login(body *authfx.LoginBody) (*models.PublicUser, string, error) {
	args := m.Called(body)

	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}

	return args.Get(0).(*models.PublicUser), args.String(1), args.Error(2)
}
