package models

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Role       types.UserRole   `gorm:"type:role;not null" json:"role"`
	FirstName  string           `gorm:"type:varchar(128);not null" json:"first_name"`
	MiddleName string           `gorm:"type:varchar(128)" json:"middle_name"`
	LastName   string           `gorm:"type:varchar(128)" json:"last_name"`
	Phone      string           `gorm:"type:varchar(15);not null" json:"phone"`
	Gender     types.UserGender `gorm:"type:gender;not null" json:"gender"`
	Email      string           `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password   string           `gorm:"type:varchar(60);not null" json:"password"`
	AvatarUrl  string           `gorm:"type:varchar(512)" json:"avartar_url"`
	SchoolNum  string           `gorm:"type:varchar(16)" json:"school_num"`
}

// Remove sensitive fields e.g. password
type PublicUser struct {
	ID         uuid.UUID
	Role       types.UserRole
	FirstName  string
	MiddleName string
	LastName   string
	Phone      string
	Gender     types.UserGender
	Email      string
	AvatarUrl  string
	SchoolNum  string
}

func (u User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:         u.ID,
		Role:       u.Role,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Gender:     u.Gender,
		Email:      u.Email,
		AvatarUrl:  u.AvatarUrl,
		SchoolNum:  u.SchoolNum,
	}
}
