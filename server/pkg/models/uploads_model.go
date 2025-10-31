package models

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/google/uuid"
)

type Upload struct {
	ObjectName string             `gorm:"type:varchar(128);primaryKey" json:"object_name"`
	UserID     uuid.UUID          `gorm:"type:uuid;primaryKey" json:"user_id"`
	Type       types.UploadType   `gorm:"type:upload_type;not null" json:"type"`
	Status     types.UploadStatus `gorm:"type:upload_status;not null" json:"status"`

	// Tells GORM that 'UserID' above refers to 'User' model
	User User `gorm:"foreignKey:UserID"`
}
