package auth_unit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/test/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// ======================== REGISTER ========================

func TestAuthController_Register_Success(t *testing.T) {
	// ------------------ Arrange ------------------
	gin.SetMode(gin.TestMode)

	mockAuthService := new(mocks.MockAuthService)
	authController := &authfx.AuthController{
		FlagConfig: &configfx.FlagConfig{
			Environment: "test",
		},
		AppConfig: &configfx.AppConfig{
			JWTExpiresIn: 24,
		},
		Logger:      zap.NewNop(),
		AuthService: mockAuthService,
	}

	registerBody := &authfx.RegisterBody{
		Role:      types.UserRoleStudent,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "johnsmith@gmail.com",
		Password:  "12345678",
		SchoolNum: "12345",
	}

	expectedUser := &models.PublicUser{
		ID:        uuid.Must(uuid.NewRandom()),
		Role:      types.UserRoleStudent,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "johnsmith@gmail.com",
		SchoolNum: "12345",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("validatedBody", registerBody)

	// Setup mock expectation
	mockAuthService.On("Register", registerBody).Return(expectedUser, "testtokenvalue", nil)

	// ------------------ Act ----------------------
	authController.Register(ctx)

	// ------------------ Assert -------------------
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Verify cookie
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, "testtokenvalue", cookies[0].Value)

	// Verify response body
	assert.Contains(t, responseBody, "user")
	userMap, _ := responseBody["user"].(map[string]any)
	assert.Equal(t, "johnsmith@gmail.com", userMap["email"])

	// Verify mock was called
	mockAuthService.AssertExpectations(t)
}

func TestAuthController_Register_MissingRequiredSchoolNum(t *testing.T) {
	testCases := []struct {
		Role types.UserRole
	}{
		{types.UserRoleStudent},
		{types.UserRoleTeacher},
	}

	for _, tc := range testCases {
		// ------------------ Arrange ------------------
		authController := &authfx.AuthController{
			Logger: zap.NewNop(),
		}

		registerBody := &authfx.RegisterBody{
			Role:      tc.Role,
			FirstName: "John",
			LastName:  "Smith",
			Phone:     "+66912345678",
			Gender:    types.UserGenderMale,
			Email:     "johnsmith@gmail.com",
			Password:  "12345678",
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("validatedBody", registerBody)

		// ------------------ Act ----------------------
		authController.Register(ctx)

		// ------------------ Assert -------------------
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Parse response body
		var responseBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody, "error")
		assert.Contains(t, responseBody["error"],
			fmt.Sprintf("%s must provide school_num", tc.Role))
	}
}

func TestAuthController_Register_UnexpectedSchoolNum(t *testing.T) {
	// ------------------ Arrange ------------------
	authController := &authfx.AuthController{
		Logger: zap.NewNop(),
	}

	registerBody := &authfx.RegisterBody{
		Role:      types.UserRoleGuardian,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "johnsmith@gmail.com",
		Password:  "12345678",
		SchoolNum: "92839",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("validatedBody", registerBody)

	// ------------------ Act ----------------------
	authController.Register(ctx)

	// ------------------ Assert -------------------
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse response body
	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "error")
	assert.Contains(t, responseBody["error"],
		fmt.Sprintf("%s should not provide school_num", types.UserRoleGuardian))
}

func TestAuthController_Register_InternalServerError(t *testing.T) {
	testCases := []struct {
		errType error
	}{
		{common.ErrDatabase},
		{common.ErrTokenGeneration},
	}

	for _, tc := range testCases {
		// ------------------ Arrange ------------------
		mockAuthService := new(mocks.MockAuthService)
		authController := &authfx.AuthController{
			Logger:      zap.NewNop(),
			AuthService: mockAuthService,
		}

		registerBody := &authfx.RegisterBody{
			Role:      types.UserRoleStudent,
			FirstName: "John",
			LastName:  "Smith",
			Phone:     "+66912345678",
			Gender:    types.UserGenderMale,
			Email:     "johnsmith@gmail.com",
			Password:  "12345678",
			SchoolNum: "1",
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("validatedBody", registerBody)

		// Setup mock expectation
		mockAuthService.On("Register", registerBody).Return(nil, "", tc.errType)

		// ------------------ Act ----------------------
		authController.Register(ctx)

		// ------------------ Assert -------------------
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse response body
		var responseBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody, "error")
		assert.Contains(t, responseBody["error"], tc.errType.Error())

		// Verify mock was called
		mockAuthService.AssertExpectations(t)
	}
}

// ======================== LOGIN ========================
func TestAuthController_Login_Success(t *testing.T) {
	// ------------------ Arrange ------------------
	mockAuthService := new(mocks.MockAuthService)
	authController := &authfx.AuthController{
		FlagConfig: &configfx.FlagConfig{
			Environment: "test",
		},
		AppConfig: &configfx.AppConfig{
			JWTExpiresIn: 24,
		},
		Logger:      zap.NewNop(),
		AuthService: mockAuthService,
	}

	loginBody := &authfx.LoginBody{
		Email:    "johnsmith@gmail.com",
		Password: "12345678",
	}

	expectedUser := &models.PublicUser{
		ID:        uuid.Must(uuid.NewRandom()),
		Role:      types.UserRoleStudent,
		FirstName: "John",
		LastName:  "Smith",
		Phone:     "+66912345678",
		Gender:    types.UserGenderMale,
		Email:     "johnsmith@gmail.com",
		SchoolNum: "12345",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("validatedBody", loginBody)

	// Setup mock expectation
	mockAuthService.On("Login", loginBody).Return(expectedUser, "testtokenvalue", nil)

	// ------------------ Act ----------------------
	authController.Login(ctx)

	// ------------------ Assert -------------------
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Verify cookie
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, "testtokenvalue", cookies[0].Value)

	// Verify response body
	assert.Contains(t, responseBody, "user")
	userMap, _ := responseBody["user"].(map[string]any)
	assert.Equal(t, "johnsmith@gmail.com", userMap["email"])

	// Verify mock was called
	mockAuthService.AssertExpectations(t)
}

func TestAuthController_Login_InternalServerError(t *testing.T) {
	testCases := []struct {
		errType error
	}{
		{common.ErrDatabase},
		{common.ErrTokenGeneration},
	}

	for _, tc := range testCases {
		// ------------------ Arrange ------------------
		mockAuthService := new(mocks.MockAuthService)
		authController := &authfx.AuthController{
			Logger:      zap.NewNop(),
			AuthService: mockAuthService,
		}

		loginBody := &authfx.LoginBody{
			Email:    "johnsmith@gmail.com",
			Password: "12345678",
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("validatedBody", loginBody)

		// Setup mock expectation
		mockAuthService.On("Login", loginBody).Return(nil, "", tc.errType)

		// ------------------ Act ----------------------
		authController.Login(ctx)

		// ------------------ Assert -------------------
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse response body
		var responseBody map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Contains(t, responseBody, "error")
		assert.Contains(t, responseBody["error"], tc.errType.Error())

		// Verify mock was called
		mockAuthService.AssertExpectations(t)
	}
}
