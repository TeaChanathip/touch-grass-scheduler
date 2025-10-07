package auth_unit_test

import (
	"testing"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ======================== REGISTER ========================

func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	mockUserService := new(mocks.MockUserService)
	authService := &authfx.AuthService{
		UserService: mockUserService,
		AppConfig: &configfx.AppConfig{
			JWTSecret:    "test-secret",
			JWTExpiresIn: 24,
		},
	}

	registerBody := authfx.RegisterBody{
		Role:      types.UserRoleStudent,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "johnsmith@gmail.com",
		Password:  "12345678",
		SchoolNum: "1",
	}

	// Setup mock expectation
	mockUserService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	user, token, err := authService.Register(registerBody)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, types.UserRoleStudent, user.Role)
	assert.Equal(t, "John", user.FirstName)
	assert.Empty(t, user.MiddleName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, "+66912345678", user.Phone)
	assert.Equal(t, types.UserGenderMale, user.Gender)
	assert.Equal(t, "johnsmith@gmail.com", user.Email)
	assert.Equal(t, "1", user.SchoolNum)

	// Verify mock was called
	mockUserService.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	mockUserService := new(mocks.MockUserService)
	authService := &authfx.AuthService{UserService: mockUserService}

	registerBody := authfx.RegisterBody{
		Role:      types.UserRoleStudent,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "duplicate@gmail.com",
		Password:  "12345678",
		SchoolNum: "1",
	}

	// Setup mock expectation
	mockUserService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(common.ErrDuplicatedEmail)

	// Act
	user, token, err := authService.Register(registerBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, common.ErrDuplicatedEmail, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Verify the mock was called exactly once
	mockUserService.AssertExpectations(t)
}

// ======================== LOGIN ========================

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockUserService := new(mocks.MockUserService)
	authService := &authfx.AuthService{
		UserService: mockUserService,
		AppConfig: &configfx.AppConfig{
			JWTSecret:    "test-secret",
			JWTExpiresIn: 24,
		},
	}

	loginBody := authfx.LoginBody{
		Email:    "johnsmith@gmail.com",
		Password: "12345678",
	}

	expectedUser := &models.User{
		Email:    "johnsmith@gmail.com",
		Password: "$2a$12$20IzYYMVPI2I79ceTEXx6upUNULaygvivZzZyBWIHb0lzJPR8P3iy", // bcrypt hash
	}

	// Setup mock expectation
	mockUserService.On("GetUserByEmail", "johnsmith@gmail.com").Return(expectedUser, nil)

	// Act
	user, token, err := authService.Login(loginBody)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "johnsmith@gmail.com", user.Email)

	// Verify the mock was called exactly once
	mockUserService.AssertExpectations(t)
}

func TestAuthService_Login_EmailNotExist(t *testing.T) {
	// Arrange
	mockUserService := new(mocks.MockUserService)
	authService := &authfx.AuthService{UserService: mockUserService}

	loginBody := authfx.LoginBody{
		Email:    "johnsmith@gmail.com",
		Password: "12345678",
	}

	// Setup mock expectation
	mockUserService.On("GetUserByEmail", "johnsmith@gmail.com").Return(nil, common.ErrUserNotFound)

	// Act
	user, token, err := authService.Login(loginBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, common.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Verify the mock was called exactly once
	mockUserService.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	mockUserService := new(mocks.MockUserService)
	authService := &authfx.AuthService{UserService: mockUserService}

	loginBody := authfx.LoginBody{
		Email:    "johnsmith@gmail.com",
		Password: "12345678",
	}

	expectedUser := &models.User{
		Email:    "johnsmith@gmail.com",
		Password: "no_matched_pwd",
	}

	// Setup mock expectation
	mockUserService.On("GetUserByEmail", "johnsmith@gmail.com").Return(expectedUser, nil)

	// Act
	user, token, err := authService.Login(loginBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, common.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Verify the mock was called exactly once
	mockUserService.AssertExpectations(t)
}
