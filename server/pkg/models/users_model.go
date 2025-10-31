package models

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/google/uuid"
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
	AvatarURL  *string          `gorm:"type:varchar(512);null;default:null" json:"avartar_url"`
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
	AvatarURL  *string          `json:"avartar_url"`
	SchoolNum  *string          `json:"school_num"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:         u.ID,
		Role:       u.Role,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Gender:     u.Gender,
		Email:      u.Email,
		AvatarURL:  u.AvatarURL,
		SchoolNum:  u.SchoolNum,
	}
}
