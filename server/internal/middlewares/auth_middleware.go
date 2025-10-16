package middlewarefx

import (
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

func (m *AuthMiddleware) Handler(role types.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := strings.TrimPrefix(ctx.Request.Header.Get("token"), "Bearer")

		if tokenString == "" {
			m.Logger.Debug("Missing token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				m.Logger.Debug(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(m.AppConfig.JWTSecret), nil
		})

		if err != nil {
			m.Logger.Debug("Invalid token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Validate token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			m.Logger.Debug("Invalid token claims")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// Extract user ID and role from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			m.Logger.Debug("Invalid user ID in token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			ctx.Abort()
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			m.Logger.Debug("Invalid role in token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			ctx.Abort()
			return
		}

		// Check if user has required role
		if role != "" && types.UserRole(userRole) != role {
			m.Logger.Debug("Insufficient permissions")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			ctx.Abort()
			return
		}

		// Set user info in context for later use
		ctx.Set("user_id", userID)
		ctx.Set("role", userRole)

		ctx.Next()
	}
}
