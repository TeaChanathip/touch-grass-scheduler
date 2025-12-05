package middlewarefx

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
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
	accessTokenString, err := ctx.Cookie("accessToken")
	if err != nil {
		return "", "", fmt.Errorf("failed retrieving access token: %w", err)
	}
	if accessTokenString == "" {
		return "", "", errors.New("access token is empty")
	}

	accessToken, err := common.ParseJWTToken(accessTokenString, m.AppConfig.JWTSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed parsing access token: %w", err)
	}

	if !accessToken.Valid {
		return "", "", errors.New("access token is not valid")
	}

	// Validate accessToken claims
	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid or missing claims")
	}

	// Extract user ID and role from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("invalid or missing user_id in claims")
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid or missing role in claims")
	}

	return userID, types.UserRole(userRole), nil
}

func (m *AuthMiddleware) HandlerWithRole(roles ...types.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, userRole, err := m.HandlerCoreLogic(ctx)
		if err != nil {
			m.Logger.Debug("Handle access token failed", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing access token"})
			return
		}

		// Check if user has required role
		if !slices.Contains(roles, userRole) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
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
			m.Logger.Debug("Handle access token failed", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing access token"})
			return
		}

		// Set user info in context for later use
		ctx.Set("user_id", userID)
		ctx.Set("role", userRole)

		ctx.Next()
	}
}
