package endpoints

import "github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"

type AuthEndpoint types.BaseStringEnum

const (
	GetRegistrationMailV1 AuthEndpoint = "api/v1/auth/registration-mail"
	RegisterV1            AuthEndpoint = "api/v1/auth/register"
	LoginV1               AuthEndpoint = "api/v1/auth/login"
	LogoutV1              AuthEndpoint = "api/v1/auth/logout"
)
