package models

import (
	"time"

	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/google/uuid"
)

type PendingUpload struct {
	ObjectKey string           `gorm:"type:varchar(128);primaryKey" json:"object_key"`
	UserID    uuid.UUID        `gorm:"type:uuid;primaryKey" json:"user_id"`
	Type      types.UploadType `gorm:"type:upload_type;not null" json:"type"`
	ExpireAt  time.Time        `gorm:"type:timestamptz;not null" json:"expire_at"`

	// Tells GORM that 'UserID' above refers to 'User' model
	User User `gorm:"foreignKey:UserID"`
}

func (PendingUpload) TableName() string {
	return "pending_uploads"
}
