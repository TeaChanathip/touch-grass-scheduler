package endpoints

import "github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"

type UsersEndpoint types.BaseStringEnum

const (
	GetMeV1                    UsersEndpoint = "api/v1/users/me"
	GetUserByIDV1              UsersEndpoint = "api/v1/users"
	UpdateUserByIDV1           UsersEndpoint = "api/v1/users"
	GetUploadAvatarSignedURLV1 UsersEndpoint = "api/v1/users/avatar-signed-url"
	HandleAvatarUploadV1       UsersEndpoint = "api/v1/users/avatar"
)
