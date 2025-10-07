package auth_unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/test/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

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

	// Create valid request body
	requestBody := map[string]any{
		"role":       "student",
		"first_name": "John",
		"last_name":  "Smith",
		"phone":      "+66912345678",
		"gender":     "male",
		"email":      "johnsmith@gmail.com",
		"password":   "12345678",
		"school_num": "12345",
	}
	jsonBody, _ := json.Marshal(requestBody)

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

	// Create proper HTTP request with body
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Setup mock expectation
	mockAuthService.On("Register", mock.AnythingOfType("authfx.RegisterBody")).Return(expectedUser, "testtokenvalue", nil)

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
	assert.Equal(t, "johnsmith@gmail.com", userMap["Email"])

	// Verify mock was called
	mockAuthService.AssertExpectations(t)
}
