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
	"golang.org/x/sync/errgroup"
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
	GetUploadAvatarSignedURL(userID uuid.UUID) (*GetUploadAvatarSignedURLResponse, error)
	UpdateUserPwdByEmail(email, newPassword string) error
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	HandleAvatarUpload(ctx context.Context, userID uuid.UUID) (*url.URL, error)
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

	publicUser, err := user.ToPublic(
		service.Logger,
		service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		return nil, err
	}

	return publicUser, nil
}

func (service *UserService) UpdateUserByID(
	userID uuid.UUID,
	body *UpdateUserBody,
) (*models.PublicUser, error) {
	var updatedUser *models.User

	// NOTE: Gorm doen't support update and return in one operation
	// Utilize transaction for atomicity
	err := service.DB.Transaction(
		func(tx *gorm.DB) error {
			// 1. Perform update
			result := tx.Model(&models.User{}).
				Where("id = ?", userID).
				Updates(&body)

			if result.Error != nil {
				service.Logger.Error("User database update failed",
					zap.String("user_id", userID.String()),
					zap.Error(result.Error),
				)
				return common.ErrDatabase
			}

			// No row affected (no user found)
			if result.RowsAffected == 0 {
				service.Logger.Debug(
					"User database update skipped",
					zap.String("reason", "user_not_found"),
					zap.String("user_id", userID.String()),
				)
				return common.ErrUserNotFound
			}

			// 2. Get updated user
			err := tx.First(&updatedUser, "id = ?", userID).Error
			if err != nil {
				service.Logger.Error("User database retrieval failed",
					zap.String("user_id", userID.String()),
					zap.Error(err),
				)
				return common.ErrDatabase
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	publicUser, err := updatedUser.ToPublic(
		service.Logger,
		service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		return nil, err
	}

	return publicUser, nil
}

func (service *UserService) GetUploadAvatarSignedURL(
	userID uuid.UUID,
) (*GetUploadAvatarSignedURLResponse, error) {
	var result *gorm.DB

	// Delete old pending uploads (if exists)
	result = service.DB.Where("user_id = ? AND type = 'avatar'", userID).
		Delete(models.PendingUpload{})
	if result.Error != nil {
		service.Logger.Error(
			"Pending upload database deletion failed",
			zap.String("user_id", userID.String()),
			zap.String("type", "pending_user_avatar"),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	// Generate ID that will be used in object name
	objectID, err := uuid.NewRandom()
	if err != nil {
		service.Logger.Error("UUID generation failed", zap.Error(err))
		return nil, common.ErrUUIDGeneration
	}

	// Generate URL
	objectKey := fmt.Sprintf("avatars/%s.webp", objectID.String())
	pendingObjectKey := fmt.Sprintf("pending/%s", objectKey)
	url, formData, err := service.generateAvatarUploadURL(pendingObjectKey)
	if err != nil {
		service.Logger.Error(
			"Signed POST URL storage generation failed",
			zap.String("user_id", userID.String()),
			zap.String("type", "pending_user_avatar"),
			zap.Error(err),
		)
		return nil, common.ErrStorage
	}

	// Create Upload entity in DB
	pendingUpload := &models.PendingUpload{
		ObjectKey: objectKey,
		UserID:    userID,
		Type:      types.UploadTypeAvatar,
		ExpireAt:  time.Now().Add(time.Hour * 24),
	}

	result = service.DB.Create(pendingUpload)
	if result.Error != nil {
		service.Logger.Error("Pending upload database creation failed",
			zap.String("user_id", userID.String()),
			zap.String("type", "pending_user_avatar"),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	response := &GetUploadAvatarSignedURLResponse{
		URL:      url.String(),
		FormData: formData,
	}

	return response, nil
}

func (service *UserService) HandleAvatarUpload(
	ctx context.Context,
	userID uuid.UUID,
) (*url.URL, error) {
	// 1. Query pending upload of user's avatar
	var pendingUpload *models.PendingUpload
	result := service.DB.Where("user_id = ? AND type = 'avatar'", userID).First(&pendingUpload)
	if result.Error != nil {
		service.Logger.Error(
			"Pending upload database retrieval failed",
			zap.String("user_id", userID.String()),
			zap.String("type", "pending_user_avatar"),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	// 2. Check if the object actually exists on the Storage
	pendingKey := fmt.Sprintf("pending/%s", pendingUpload.ObjectKey)
	_, err := service.StorageClient.StatObject(ctx,
		service.AppConfig.StorageBucketName,
		pendingKey,
		minio.StatObjectOptions{})
	if err != nil {
		var minioErr minio.ErrorResponse
		if errors.As(err, &minioErr) && minioErr.Code == "NoSuchKey" {
			service.Logger.Debug(
				"Object storage upload skipped",
				zap.String("reason", "object_not_found"),
				zap.String("type", "pending_user_avatar"),
				zap.String("", pendingUpload.ObjectKey),
			)
			return nil, common.ErrStorageObjectNotFound
		}
		service.Logger.Error(
			"Object info storage retrieval failed",
			zap.String("object_key", pendingKey),
			zap.String("type", "pending_user_avatar"),
			zap.Error(err),
		)
		return nil, common.ErrStorage
	}

	var updatedUser *models.User
	var oldAvatarKey *string
	err = service.DB.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB

		// 3. Get user by ID
		result = tx.Where("id = ?", userID).First(&updatedUser)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User not found
			service.Logger.Debug(
				"User avatar database upload skipped",
				zap.String("reason", "user_not_found"),
				zap.String("user_id", userID.String()),
			)
			return common.ErrUserNotFound
		} else if result.Error != nil {
			// Other errors
			service.Logger.Error("User database retrieval failed",
				zap.String("user_id", userID.String()),
				zap.Error(result.Error),
			)
			return common.ErrDatabase
		}

		oldAvatarKey = updatedUser.AvatarKey

		// Error group for the following database operation go routines
		gDB, errGroupDBCtx := errgroup.WithContext(ctx)

		// 4. Update 'avatar_url' of user in DB
		gDB.Go(func() error {
			result = tx.WithContext(errGroupDBCtx).Model(&updatedUser).
				Update("avatar_key", &pendingUpload.ObjectKey)
			if result.Error != nil {
				// Other errors
				service.Logger.Error(
					"User avatar database update failed",
					zap.String("user_id", userID.String()),
					zap.Error(result.Error))
				return common.ErrDatabase
			}
			return nil
		})

		// 5. Delete pending upload in Database
		gDB.Go(func() error {
			result = tx.
				WithContext(errGroupDBCtx).
				Where("user_id = ? AND object_key = ? AND type = 'avatar'", userID, pendingUpload.ObjectKey).
				Delete(models.PendingUpload{})
			if result.Error != nil {
				service.Logger.Error(
					"Pending upload database deletion failed",
					zap.String("user_id", userID.String()),
					zap.String("type", "pending_user_avatar"),
					zap.Error(result.Error),
				)
				return common.ErrDatabase
			}
			if result.RowsAffected == 0 {
				// "No pending upload deleted",
				service.Logger.Error(
					"Pending upload database deletion skipped",
					zap.String("reason", "pending_upload_not_found"),
					zap.String("type", "pending_user_avatar"),
					zap.String("user_id", userID.String()),
				)
				return common.ErrPendingUploadNotFound
			}
			return nil
		})

		// Return error if any go routines error
		return gDB.Wait()
	})
	if err != nil {
		return nil, err
	}

	// 6. Copy and rename pending object in Storage
	src := minio.CopySrcOptions{
		Bucket: service.AppConfig.StorageBucketName,
		Object: pendingKey,
	}
	dst := minio.CopyDestOptions{
		Bucket: service.AppConfig.StorageBucketName,
		Object: pendingUpload.ObjectKey,
	}
	_, err = service.StorageClient.CopyObject(ctx, dst, src)
	if err != nil {
		service.Logger.Error(
			"Object storage copy failed",
			zap.String("from_object_key", pendingKey),
			zap.String("to_object_key", pendingUpload.ObjectKey),
			zap.String("type", "user_avatar"),
			zap.Error(err),
		)
		return nil, common.ErrStorage
	}

	// Error group for the following storage operation go routines
	gStorage, errGroupStorageCtx := errgroup.WithContext(ctx)

	// 7. Delete old avatar from Storage (if exists)
	if oldAvatarKey != nil {
		gStorage.Go(func() error {
			err := service.StorageClient.RemoveObject(errGroupStorageCtx,
				service.AppConfig.StorageBucketName,
				*oldAvatarKey,
				minio.RemoveObjectOptions{})
			if err != nil {
				service.Logger.Error(
					"Object storage deletion failed",
					zap.String("object_key", *oldAvatarKey),
					zap.String("type", "user_avatar"),
					zap.Error(err),
				)
				return common.ErrStorage
			}
			return nil
		})
	}

	// 8. Delete the pending object in Storage
	gStorage.Go(func() error {
		err = service.StorageClient.RemoveObject(errGroupStorageCtx,
			service.AppConfig.StorageBucketName,
			pendingKey,
			minio.RemoveObjectOptions{})
		if err != nil {
			service.Logger.Error(
				"Object storage deletion failed",
				zap.String("object_key", pendingKey),
				zap.String("type", "pending_user_avatar"),
				zap.Error(err),
			)
			return common.ErrStorage
		}
		return nil
	})

	if err := gStorage.Wait(); err != nil {
		return nil, err
	}

	signedURL, err := service.StorageClient.PresignedGetObject(ctx,
		service.AppConfig.StorageBucketName,
		pendingUpload.ObjectKey,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn),
		nil)
	if err != nil {
		service.Logger.Error(
			"Signed GET URL storage generation failed",
			zap.String("type", "user_avatar"),
			zap.Error(err),
		)
		return nil, common.ErrStorage
	}

	return signedURL, nil
}

// ======================== HELPER METHODS ========================

func (service *UserService) CreateUser(user *models.User) error {
	// Hash the password
	hashed, err := common.HashPassword(user.Password)
	if err != nil {
		service.Logger.Error("Password hashing failed", zap.Error(err))
		return common.ErrPasswordHashing
	}

	// Replace password with hashed
	user.Password = hashed

	// Create new User in DB
	result := service.DB.Create(&user)

	if result.Error != nil {
		// Check for PostgreSQL unique constraint violation
		// I know this looks absurd, but it is the simplest solution
		if strings.Contains(result.Error.Error(), "SQLSTATE 23505") {
			service.Logger.Debug(
				"User database creation skipped",
				zap.String("reason", "email_duplicated"),
				zap.String("email", user.Email),
			)
			return common.ErrDuplicatedEmail
		}

		// Other database errors
		service.Logger.Error(
			"User database creation failed",
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
		service.Logger.Debug(
			"User database retrieval skipped",
			zap.String("reason", "user_not_found"),
			zap.String("email", email),
		)
		return nil, common.ErrUserNotFound
	} else if result.Error != nil {
		// Other errors
		service.Logger.Error(
			// "Database error while fetching user with email",
			"User database retrieval failed",
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
		service.Logger.Debug(
			"User database retrieval skipped",
			zap.String("reason", "user_id_not_found"),
			zap.String("user_id", userID.String()),
		)
		return nil, common.ErrUserNotFound
	} else if result.Error != nil {
		// Other errors
		service.Logger.Error(
			"User database retrieval failed",
			zap.String("id", userID.String()),
			zap.Error(result.Error),
		)
		return nil, common.ErrDatabase
	}

	return user, nil
}

func (service *UserService) UpdateUserPwdByEmail(email, newPassword string) error {
	// Hash the password
	hashed, err := common.HashPassword(newPassword)
	if err != nil {
		service.Logger.Error("Password hashing failed", zap.String("email", email), zap.Error(err))
		return common.ErrPasswordHashing
	}

	err = service.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("email = ?", email).Update("password", hashed)
		// Must be only one user that affected
		if result.Error != nil || result.RowsAffected != 1 {
			service.Logger.Error(
				"User password database update failed",
				zap.String("email", email),
				zap.Error(result.Error),
			)
			return common.ErrDatabase
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) generateAvatarUploadURL(
	objectKey string,
) (*url.URL, map[string]string, error) {
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
