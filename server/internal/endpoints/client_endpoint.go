package endpoints

import "github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"

type ClientEndpoint types.BaseStringEnum

const (
	ClientRegistrationVerification ClientEndpoint = "register" // token is required
	ClientForgotPwd                ClientEndpoint = "forgot-password"
	ClientResetPwd                 ClientEndpoint = "reset-password" // token is required
)
