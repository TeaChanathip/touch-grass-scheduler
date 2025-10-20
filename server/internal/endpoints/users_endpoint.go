package endpoints

import "github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"

type UsersEndpoint types.BaseStringEnum

const (
	GetMeV1         UsersEndpoint = "api/v1/users/me"
	GetUserWithIDV1 UsersEndpoint = "api/v1/users"
)
