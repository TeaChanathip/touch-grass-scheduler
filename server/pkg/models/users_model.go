package models

import (
	"context"
	"fmt"
	"time"

	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type User struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Role       types.UserRole   `gorm:"type:role;not null" json:"role"`
	FirstName  string           `gorm:"type:varchar(128);not null" json:"first_name"`
	MiddleName string           `gorm:"type:varchar(128);null;default:''" json:"middle_name"`
	LastName   string           `gorm:"type:varchar(128);null;default:''" json:"last_name"`
	Phone      string           `gorm:"type:varchar(15);not null" json:"phone"`
	Gender     types.UserGender `gorm:"type:gender;not null" json:"gender"`
	Email      string           `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password   string           `gorm:"type:varchar(60);not null" json:"password"`
	AvatarKey  *string          `gorm:"type:varchar(512);null;default:null" json:"avatar_key"`
	SchoolNum  *string          `gorm:"type:varchar(16);null;default:null" json:"school_num"`
}

// PublicUser Remove sensitive fields e.g. password
type PublicUser struct {
	ID         uuid.UUID        `json:"id"`
	Role       types.UserRole   `json:"role"`
	FirstName  string           `json:"first_name"`
	MiddleName string           `json:"middle_name"`
	LastName   string           `json:"last_name"`
	Phone      string           `json:"phone"`
	Gender     types.UserGender `json:"gender"`
	Email      string           `json:"email"`
	AvartarURL *string          `json:"avatar_url"`
	SchoolNum  *string          `json:"school_num"`
}

func (u *User) ToPublic(storageClient *minio.Client, bucketName string, expires time.Duration) (*PublicUser, error) {
	var avatarURL *string = nil

	if u.AvatarKey != nil {
		ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancle()

		signedURL, err := storageClient.PresignedGetObject(ctx, bucketName, *u.AvatarKey, expires, nil)
		if err != nil {
			return nil, fmt.Errorf("failed getting presigned url: %w", err)
		}

		signedURLStr := signedURL.String()
		avatarURL = &signedURLStr
	}

	publicUser := &PublicUser{
		ID:         u.ID,
		Role:       u.Role,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Gender:     u.Gender,
		Email:      u.Email,
		AvartarURL: avatarURL,
		SchoolNum:  u.SchoolNum,
	}

	return publicUser, nil
}
