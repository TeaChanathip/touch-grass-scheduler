package authfx

import (
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	usersfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/users"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthServiceParams struct {
	fx.In
	AppConfig   *configfx.AppConfig
	Logger      *zap.Logger
	UserService usersfx.UserServiceInterface
}

type AuthService struct {
	AppConfig   *configfx.AppConfig
	Logger      *zap.Logger
	UserService usersfx.UserServiceInterface
}

func NewAuthService(params AuthServiceParams) *AuthService {
	return &AuthService{
		AppConfig:   params.AppConfig,
		Logger:      params.Logger,
		UserService: params.UserService,
	}
}

// ======================== METHODS ========================

func (service *AuthService) Register(registerBody RegisterBody) (*models.PublicUser, string, error) {
	// TODO: Add logic to check if SchoolNumber is valid
	// TODO: Send the verification link to the user's email

	// Create new user
	user := registerBody.ToUserModel()
	if err := service.UserService.CreateUser(user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user.ToPublic(), token, nil
}

func (service *AuthService) Login(loginBody LoginBody) (*models.PublicUser, string, error) {
	user, err := service.UserService.GetUserByEmail(loginBody.Email)
	if err != nil {
		return nil, "", common.ErrInvalidCredentials
	}

	// Compare password with hashed
	if !common.CheckHashedPassword(loginBody.Password, user.Password) {
		return nil, "", common.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user.ToPublic(), token, nil
}

// ======================== HELPER METHODS ========================

func (service *AuthService) generateToken(user *models.User) (string, error) {
	// Get expired duration from ENV
	exp := time.Now().Add(time.Duration(service.AppConfig.JWTExpiresIn) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     jwt.NewNumericDate(exp),
		})

	singedToken, err := token.SignedString([]byte(service.AppConfig.JWTSecret))
	if err != nil {
		service.Logger.Error("Internal error while signing the JWT:", zap.Error(err))
		return "", common.ErrTokenGeneration
	}

	return singedToken, nil
}
