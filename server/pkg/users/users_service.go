package usersfx

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserServiceParams struct {
	fx.In
	AppConfig     *configfx.AppConfig
	Logger        *zap.Logger
	DB            *gorm.DB
	StorageClient *minio.Client
}

type UserService struct {
	AppConfig     *configfx.AppConfig
	Logger        *zap.Logger
	DB            *gorm.DB
	StorageClient *minio.Client
}

type UserServiceInterface interface {
	GetPublicUserByID(userID uuid.UUID) (*models.PublicUser, error)
	UpdateUserByID(userID uuid.UUID, body *UpdateUserBody) (*models.PublicUser, error)
	GetUploadAvatarSignedURL(userID uuid.UUID) (map[string]any, error)
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	HandleAvatarUpload(userID uuid.UUID) (*models.PublicUser, error)
	generateAvatarUploadURL(objectName string) (*url.URL, map[string]string, error)
}

// Verify interface implementation at compile time
var _ UserServiceInterface = (*UserService)(nil)

func NewUserService(params UserServiceParams) UserServiceInterface {
	return &UserService{
		AppConfig:     params.AppConfig,
		Logger:        params.Logger,
		DB:            params.DB,
		StorageClient: params.StorageClient,
	}
}

// ======================== BUSINESS LOGIC METHODS ========================

func (service *UserService) GetPublicUserByID(userID uuid.UUID) (*models.PublicUser, error) {
	user, err := service.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	publicUser, err := user.ToPublic(service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Duration(time.Hour*time.Duration(service.AppConfig.JWTExpiresIn)))
	if err != nil {
		return nil, err
	}

	return publicUser, nil
}

func (service *UserService) UpdateUserByID(userID uuid.UUID, body *UpdateUserBody) (*models.PublicUser, error) {
	var updatedUser *models.User

	// NOTE: Gorm doen't support update and return in one operation
	// Utilize transaction for atomicity
	err := service.DB.Transaction(
		func(tx *gorm.DB) error {
			// Operation1: Perform update
			result := tx.Model(&models.User{}).
				Where("id = ?", userID).
				Updates(&body)

			if result.Error != nil {
				service.Logger.Error("Database error while updating user",
					zap.String("id", userID.String()),
					zap.Error(result.Error),
				)
				return common.ErrDatabase
			}

			// No row affected (no user found)
			if result.RowsAffected == 0 {
				service.Logger.Debug("User not found for update", zap.String("id", userID.String()))
				return common.ErrUserNotFound
			}

			// Operation2: Get updated user
			err := tx.First(&updatedUser, "id = ?", userID).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					service.Logger.Error("Updated user but could not fetch it (race condition?)",
						zap.String("id", userID.String()),
					)
					return common.ErrUserNotFound
				}

				service.Logger.Error("Database error while fetching updated user",
					zap.String("id", userID.String()),
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

	publicUser, err := updatedUser.ToPublic(service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Duration(time.Hour*time.Duration(service.AppConfig.JWTExpiresIn)))
	if err != nil {
		service.Logger.Error("Error getting signed URL for user's avatar", zap.Error(err))
		return nil, err
	}

	return publicUser, nil
}

func (service *UserService) GetUploadAvatarSignedURL(userID uuid.UUID) (map[string]any, error) {
	// Generate ID that will be used in object name
	objectID, err := uuid.NewRandom()
	if err != nil {
		service.Logger.Error("Error while generating UUID", zap.Error(err))
		return nil, common.ErrUUIDGenerating
	}

	// Generate URL
	objectKey := fmt.Sprintf("avatars/%s.webp", objectID.String())
	pendingObjectKey := fmt.Sprintf("pending/%s", objectKey)
	url, formData, err := service.generateAvatarUploadURL(pendingObjectKey)
	if err != nil {
		service.Logger.Error("", zap.Error(err))
		return nil, common.ErrStorage
	}

	// Create Upload entity in DB
	upload := &models.PendingUpload{
		ObjectKey: objectKey,
		UserID:    userID,
		Type:      types.UploadTypeAvatar,
	}

	result := service.DB.Create(upload)
	if result.Error != nil {
		service.Logger.Error("Database error while creating user",
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	response := map[string]any{
		"url":        url.String(),
		"form_data":  formData,
		"object_key": objectKey,
	}

	return response, nil
}

func (service *UserService) HandleAvatarUpload(userID uuid.UUID) (*models.PublicUser, error) {
	// Query pending upload of user's avatar
	var pendingUpload *models.PendingUpload
	result := service.DB.Where("user_id = ? AND type = 'avatar'", userID.String()).First(&pendingUpload)
	if result.Error != nil {
		service.Logger.Error("Database error while getting pending upload with UserID", zap.Error(result.Error))
		return nil, common.ErrDatabase
	}

	// Check if the object actually exists on the Storage
	ctx := context.Background()
	_, err := service.StorageClient.StatObject(ctx,
		service.AppConfig.StorageBucketName,
		fmt.Sprintf("pending/%s", pendingUpload.ObjectKey),
		minio.StatObjectOptions{})
	if err != nil {
		minioErr, ok := err.(minio.ErrorResponse)
		if ok && minioErr.Code == "NoSuchKey" {
			service.Logger.Sugar().Debugf("%s not found in Storage", pendingUpload.ObjectKey)
			return nil, common.ErrStorageObjectNotFound
		}
		service.Logger.Error("Storage error while getting ObjectInfo", zap.Error(err))
		return nil, common.ErrStorage
	}

	var updatedUser *models.User
	err = service.DB.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB

		// Get user by ID
		result = tx.Where("id = ?", userID).First(&updatedUser)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User not found
			service.Logger.Debug("User not found", zap.String("id", userID.String()))
			return common.ErrUserNotFound
		} else if result.Error != nil {
			// Other errors
			service.Logger.Error("Database error while getting user with id",
				zap.String("id", userID.String()),
				zap.Error(result.Error),
			)
			return common.ErrDatabase
		}

		oldAvatarKey := updatedUser.AvatarKey

		// Update 'avatar_url' of user in DB
		result = tx.Model(&updatedUser).
			Update("avatar_key", &pendingUpload.ObjectKey)
		if result.Error != nil {
			// Other errors
			service.Logger.Error("Database error while updating user's avatar_key",
				zap.String("id", userID.String()),
				zap.Error(result.Error))
			return common.ErrDatabase
		}

		// Delete pending upload
		result = tx.
			Where("user_id = ? AND object_key = ? AND type = 'avatar'", userID, pendingUpload.ObjectKey).
			Delete(models.PendingUpload{})
		if result.Error != nil {
			service.Logger.Error("Database error while deleting pending upload",
				zap.Error(result.Error),
			)
			return common.ErrDatabase
		}
		if result.RowsAffected == 0 {
			service.Logger.Warn("No pending upload deleted",
				zap.String("user_id", userID.String()),
				zap.String("object_key", pendingUpload.ObjectKey))
			return common.ErrPendingUploadNotFound
		}

		// Delete old avatar from Storage (if exists)
		if oldAvatarKey != nil {
			ctx := context.Background()
			err := service.StorageClient.RemoveObject(ctx,
				service.AppConfig.StorageBucketName,
				*oldAvatarKey,
				minio.RemoveObjectOptions{})
			if err != nil {
				service.Logger.Error("Storage error while deleting old avatar", zap.Error(err))
				return common.ErrStorage
			}
		}

		// Remove "pending/" prefix from the object in Storage
		src := minio.CopySrcOptions{
			Bucket: service.AppConfig.StorageBucketName,
			Object: fmt.Sprintf("pending/%s", pendingUpload.ObjectKey),
		}
		dst := minio.CopyDestOptions{
			Bucket: service.AppConfig.StorageBucketName,
			Object: pendingUpload.ObjectKey,
		}
		_, err = service.StorageClient.CopyObject(ctx, dst, src)
		if err != nil {
			service.Logger.Error("Storage error while moving avatar from pending to avatars", zap.Error(err))
			return common.ErrStorage
		}

		// Delete the pending object
		err = service.StorageClient.RemoveObject(ctx,
			service.AppConfig.StorageBucketName,
			fmt.Sprintf("pending/%s", pendingUpload.ObjectKey),
			minio.RemoveObjectOptions{})
		if err != nil {
			service.Logger.Error("Storage error while deleting pending avatar object", zap.Error(err))
			return common.ErrStorage
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	publicUser, err := updatedUser.ToPublic(service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Duration(time.Hour*time.Duration(service.AppConfig.JWTExpiresIn)))
	if err != nil {
		service.Logger.Error("Error getting signed URL for user's avatar", zap.Error(err))
		return nil, err
	}

	return publicUser, nil
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

func (service *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user *models.User

	result := service.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User not found
		service.Logger.Debug("User not found", zap.String("id", userID.String()))
		return nil, common.ErrUserNotFound
	} else if result.Error != nil {
		// Other errors
		service.Logger.Error("Database error while getting user with id",
			zap.String("id", userID.String()),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	return user, nil
}

func (service *UserService) generateAvatarUploadURL(objectKey string) (*url.URL, map[string]string, error) {
	// Create upload policy
	policy := minio.NewPostPolicy()

	// Set the bucket and object key
	if err := policy.SetBucket(service.AppConfig.StorageBucketName); err != nil {
		return nil, nil, err
	}
	if err := policy.SetKey(objectKey); err != nil {
		return nil, nil, err
	}

	// Set an expiration
	if err := policy.SetExpires(time.Now().Add(3 * time.Minute)); err != nil {
		return nil, nil, err
	}

	// Set size limit
	minLen := int64(1)
	maxLen := int64(2 * 1024 * 1024) // 2 MB
	if err := policy.SetContentLengthRange(minLen, maxLen); err != nil {
		return nil, nil, err
	}

	// Generate signed URL
	ctx := context.Background()
	url, formData, err := service.StorageClient.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}

	return url, formData, nil
}

// TODO: Avatar Upload
// 1. Add file size limit at the client
// 2. Add cropping feature at the client
// [Done] 3. Add upload policy of file size limit at the backend
// 4. (choice 1) Remove old pending upload request in DB when there's a new request
// 4. (choice 2) Add CORN job to remove old pending upload requests in DB
// 5. Add CORN job to remove pending upload objects in Storage that are too old
