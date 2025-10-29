package usersfx

import (
	"errors"
	"strings"

	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/google/uuid"
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

type UserServiceInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUserByID(id uuid.UUID, body *UpdateUserBody) (*models.User, error)
}

// Verify interface implementation at compile time
var _ UserServiceInterface = (*UserService)(nil)

func NewUserService(params UserServiceParams) UserServiceInterface {
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
	var user *models.User

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

	return user, nil
}

func (service *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user *models.User

	result := service.DB.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User not found
		service.Logger.Debug("User not found", zap.String("id", id.String()))
		return nil, common.ErrUserNotFound
	} else if result.Error != nil {
		// Other errors
		service.Logger.Error("Database error while fetching user with id",
			zap.String("id", id.String()),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	return user, nil
}

func (service *UserService) UpdateUserByID(id uuid.UUID, body *UpdateUserBody) (*models.User, error) {
	var updatedUser *models.User

	// NOTE: Gorm doen't support update and return in one operation
	// Utilize transaction for atomicity
	err := service.DB.Transaction(
		func(tx *gorm.DB) error {
			// Operation1: Perform update
			result := tx.Model(&models.User{}).
				Where("id = ?", id).
				Updates(&body)

			if result.Error != nil {
				service.Logger.Error("Database error while updating user",
					zap.String("id", id.String()),
					zap.Error(result.Error),
				)
				return common.ErrDatabase
			}

			// No row affected (no user found)
			if result.RowsAffected == 0 {
				service.Logger.Debug("User not found for update", zap.String("id", id.String()))
				return common.ErrUserNotFound
			}

			// Operation2: Get updated user
			err := tx.First(&updatedUser, "id = ?", id).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					service.Logger.Error("Updated user but could not fetch it (race condition?)",
						zap.String("id", id.String()),
					)
					return common.ErrUserNotFound
				}

				service.Logger.Error("Database error while fetching updated user",
					zap.String("id", id.String()),
					zap.Error(err),
				)
				// Return error to trigger a rollback
				return common.ErrDatabase
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
