package middlewares_unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	middlewarefx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test struct with validation tags
type TestRequestBody struct {
	Name     string `json:"name" binding:"required,min=2,max=50,alpha"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,min=1,max=120"`
	Phone    string `json:"phone" binding:"required,e164"`
	Role     string `json:"role" binding:"required,oneof=admin user guest"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func TestRequestBodyValidator_Success(t *testing.T) {
	// ------------------ Arrange ------------------
	validBody := map[string]any{
		"name":     "John",
		"email":    "john@example.com",
		"age":      25,
		"phone":    "+1234567890",
		"role":     "user",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(validBody)

	// Create test router
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	// Track if next handler was called
	var nextHandlerCalled bool
	var validatedBodyFromContext *TestRequestBody

	// Set up route with middleware and handler
	requestBodyValidator := &middlewarefx.RequestBodyValidator{Logger: zap.NewNop()}

	router.POST("/test",
		requestBodyValidator.Handler("test", TestRequestBody{}),
		func(c *gin.Context) {
			nextHandlerCalled = true

			// Verify validated body is set in context
			validatedBody, exists := c.Get("validatedBody")
			assert.True(t, exists)

			testBody, ok := validatedBody.(*TestRequestBody)
			assert.True(t, ok)
			validatedBodyFromContext = testBody

			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	// Create request
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// ------------------ Act ----------------------
	router.ServeHTTP(w, req)

	// ------------------ Assert -------------------
	assert.True(t, nextHandlerCalled)
	assert.Equal(t, 200, w.Code)

	// Verify the validated body
	assert.NotNil(t, validatedBodyFromContext)
	assert.Equal(t, "John", validatedBodyFromContext.Name)
	assert.Equal(t, "john@example.com", validatedBodyFromContext.Email)
	assert.Equal(t, 25, validatedBodyFromContext.Age)
	assert.Equal(t, "+1234567890", validatedBodyFromContext.Phone)
	assert.Equal(t, "user", validatedBodyFromContext.Role)
	assert.Equal(t, "password123", validatedBodyFromContext.Password)
}

func TestRequestBodyValidator_Error(t *testing.T) {
	// ------------------ Arrange ------------------
	validBody := map[string]any{
		"name":  "John2342",
		"email": "johnexample.com",
		"age":   20,
		"phone": "+66123456789",
		"role":  "teacher",
	}
	jsonBody, _ := json.Marshal(validBody)

	// Create test context
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	// Create middleware
	requestBodyValidator := &middlewarefx.RequestBodyValidator{Logger: zap.NewNop()}
	middleware := requestBodyValidator.Handler("test", TestRequestBody{})

	// ------------------ Act ----------------------
	middleware(ctx)

	// ------------------ Assert -------------------
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse response body
	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "error")
	assert.Contains(t, responseBody["error"], "Name")
	assert.Contains(t, responseBody["error"], "Email")
	assert.NotContains(t, responseBody["error"], "Age")
	assert.NotContains(t, responseBody["error"], "Phone")
	assert.Contains(t, responseBody["error"], "Role")
	assert.Contains(t, responseBody["error"], "Password")
}
