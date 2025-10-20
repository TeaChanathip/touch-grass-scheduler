package endpoints

import "github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"

type UsersEndpoint types.BaseStringEnum

const (
	GetUserV1 UsersEndpoint = "api/v1/users"
)
