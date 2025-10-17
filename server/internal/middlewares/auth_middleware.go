package middlewarefx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthMiddlewareParams struct {
	fx.In
	AppConfig *configfx.AppConfig
	Logger    *zap.Logger
}

type AuthMiddleware struct {
	AppConfig *configfx.AppConfig
	Logger    *zap.Logger
}

func NewAuthMiddleware(params AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		AppConfig: params.AppConfig,
		Logger:    params.Logger,
	}
}

func (m *AuthMiddleware) HandlerCoreLogic(ctx *gin.Context) (string, types.UserRole, error) {
	tokenString := strings.TrimSpace(strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer"))

	if tokenString == "" {
		return "", "", errors.New("missing token")
	}

	m.Logger.Debug(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.AppConfig.JWTSecret), nil
	})
	if err != nil {
		return "", "", errors.New("invalid token")
	}

	// Validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	// Extract user ID and role from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in token")
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid role in token")
	}

	return userID, types.UserRole(userRole), nil
}

func (m *AuthMiddleware) HandlerWithRole(role types.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, userRole, err := m.HandlerCoreLogic(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// Check if user has required role
		if role != "" && userRole != role {
			m.Logger.Debug("Insufficient permissions")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			ctx.Abort()
			return
		}

		// Set user info in context for later use
		ctx.Set("user_id", userID)
		ctx.Set("role", userRole)

		ctx.Next()
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, userRole, err := m.HandlerCoreLogic(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// Set user info in context for later use
		ctx.Set("user_id", userID)
		ctx.Set("role", userRole)

		ctx.Next()
	}
}
