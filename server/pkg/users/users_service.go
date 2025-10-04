package usersfx

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserServiceParams struct {
	fx.In
	Logger *zap.Logger
	DB     *gorm.DB
}

type UserService struct {
	Logger *zap.Logger
	DB     *gorm.DB
}

func NewUserService(params UserServiceParams) *UserService {
	return &UserService{
		Logger: params.Logger,
		DB:     params.DB,
	}
}

func (service *UserService) CreateUser(user *models.User) error {
	// Hash the password
	hashed, err := common.HashPassword(user.Password)
	if err != nil {
		service.Logger.Error("An error occured while hashing the password:", zap.Error(err))
		return err
	}

	// Replace the plain password with the hashed password
	user.Password = hashed

	// Create new User in DB
	result := service.DB.Create(&user)

	return result.Error
}

func (service *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	// Query by email
	result := service.DB.Where("email = ?", email).First(&user)

	return &user, result.Error
}
