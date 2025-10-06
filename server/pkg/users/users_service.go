package usersfx

import (
	"errors"
	"strings"

	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
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

// ======================== HELPER METHODS ========================

func (service *UserService) CreateUser(user *models.User) error {
	// Hash the password
	hashed, err := common.HashPassword(user.Password)
	if err != nil {
		service.Logger.Error("Internal error while hashing the password:", zap.Error(err))
		return common.ErrPasswordHashing
	}

	// Replace password with hashed
	user.Password = hashed

	// Create new User in DB
	result := service.DB.Create(&user)

	if result.Error != nil {
		// Check for PostgreSQL unique constraint violation
		// I know this looks absurd, but this is the simplest solution
		if strings.Contains(result.Error.Error(), "SQLSTATE 23505") {
			service.Logger.Debug("Email is duplicated", zap.String("email", user.Email))
			return common.ErrDuplicatedEmail
		}

		// Other database errors
		service.Logger.Error("Database error while creating user",
			zap.Error(result.Error),
		)
		return common.ErrDatabase
	}

	return nil
}

func (service *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	// Query user by email
	result := service.DB.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User not found
		service.Logger.Debug("User not found", zap.String("email", email))
		return nil, common.ErrUserNotFound
	} else if result.Error != nil {
		// Other errors
		service.Logger.Error("Database error while fetching user with email",
			zap.String("email", email),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	return &user, nil
}
