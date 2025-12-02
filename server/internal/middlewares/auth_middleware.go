package middlewarefx

import (
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
		return "", "", common.ErrMissingToken
	}

	accessToken, err := common.ParseJWTToken(accessTokenString, m.AppConfig.JWTSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed parsing access token: %w", err)
	}

	// Validate accessToken claims
	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", common.ErrVariableParsing
	}
	if !accessToken.Valid {
		return "", "", common.ErrInvalidCredentials
	}

	// Extract user ID and role from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		// return "", "", errors.New("invalid user ID in accessToken")
		return "", "", common.ErrMissingClaims
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed asserting string type to user role from claims: %w", err)
	}

	return userID, types.UserRole(userRole), nil
}

func (m *AuthMiddleware) HandlerWithRole(roles ...types.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, userRole, err := m.HandlerCoreLogic(ctx)
		if err != nil {
			m.Logger.Info("Error on AuthMiddleware with role", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
			m.Logger.Info("Error on AuthMiddleware", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set user info in context for later use
		ctx.Set("user_id", userID)
		ctx.Set("role", userRole)

		ctx.Next()
	}
}
