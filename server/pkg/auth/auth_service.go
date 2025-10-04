package authfx

import (
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type AuthServiceParams struct {
	fx.In
	AppConfig *configfx.AppConfig
	DB        *gorm.DB
}

type AuthService struct {
	AppConfig *configfx.AppConfig
	DB        *gorm.DB
}

func NewAuthService(params AuthServiceParams) *AuthService {
	return &AuthService{
		AppConfig: params.AppConfig,
		DB:        params.DB,
	}
}

func (service *AuthService) GenerateToken(user *models.User) (string, error) {
	// Get expired duration from ENV
	exp := time.Now().Add(time.Duration(service.AppConfig.JWTExpiresIn) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     jwt.NewNumericDate(exp),
		})

	return token.SignedString([]byte(service.AppConfig.JWTSecret))
}
