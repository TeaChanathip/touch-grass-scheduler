package authfx

import (
	"errors"
	"net/mail"
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	mailfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/mail"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	usersfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthServiceParams struct {
	fx.In
	AppConfig     *configfx.AppConfig
	Logger        *zap.Logger
	MailService   *mailfx.MailService
	UserService   usersfx.UserServiceInterface
	StorageClient *minio.Client
}

type AuthService struct {
	AppConfig     *configfx.AppConfig
	Logger        *zap.Logger
	MailService   *mailfx.MailService
	UserService   usersfx.UserServiceInterface
	StorageClient *minio.Client
}

// Verify interface implementation at compile time
var _ AuthServiceInterface = (*AuthService)(nil)

type AuthServiceInterface interface {
	GetRegistrationMail(email string) error
	Register(registrationTokenString string, body *RegisterBody) (*models.PublicUser, string, error)
	Login(body *LoginBody) (*models.PublicUser, string, error)
	GetResetPwdMail(email string) error
	ResetPwd(body *ResetPwdBody) error
}

func NewAuthService(params AuthServiceParams) AuthServiceInterface {
	return &AuthService{
		AppConfig:     params.AppConfig,
		Logger:        params.Logger,
		MailService:   params.MailService,
		UserService:   params.UserService,
		StorageClient: params.StorageClient,
	}
}

// ======================== BUSINESS LOGIC METHODS ========================

func (service *AuthService) GetRegistrationMail(email string) error {
	// Check if email already existed
	user, err := service.UserService.GetUserByEmail(email)
	if err != nil && !errors.Is(err, common.ErrUserNotFound) {
		return err
	}

	// Send warning email if user already existed
	if user != nil {
		err = service.MailService.SendRegistrationWarning(user)
		return err
	}

	// Generate registration token
	var registrationToken string
	registrationToken, err = service.generateActionToken(email,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		return err
	}

	// Send verification email if it is new user
	err = service.MailService.SendRegistrationVerification(email, registrationToken)
	return err
}

func (service *AuthService) Register(
	registrationTokenStr string,
	body *RegisterBody,
) (*models.PublicUser, string, error) {
	// TODO: Add logic to check if SchoolNumber is valid

	// Parse registerToken
	registrationToken, err := common.ParseJWTToken(
		registrationTokenStr,
		service.AppConfig.JWTSecret,
	)
	if err != nil {
		service.Logger.Debug("Registration token parse failed", zap.Error(err))
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, "", common.ErrActionTokenExpired
		}
		return nil, "", common.ErrActionTokenParsing
	}
	if !registrationToken.Valid {
		return nil, "", common.ErrInvalidActionToken
	}

	// Get email from accessToken claims
	claims, ok := registrationToken.Claims.(jwt.MapClaims)
	if !ok {
		service.Logger.Debug("Registration token type assertion failed")
		return nil, "", common.ErrActionTokenClaimsRetrieval
	}
	email, ok := claims["email"].(string)
	if !ok {
		service.Logger.Debug(
			"Registration token claims retrieval failed",
			zap.String("key", "email"),
		)
		return nil, "", common.ErrActionTokenClaimsRetrieval
	}

	// Check if email is valid
	if _, err = mail.ParseAddress(email); err != nil {
		service.Logger.Debug("Email invalid or missing", zap.String("email", email))
		return nil, "", common.ErrActionTokenClaimsRetrieval
	}

	// Create new user
	user := body.ToUserModel()
	user.Email = email
	if err := service.UserService.CreateUser(user); err != nil {
		return nil, "", err
	}
	publicUser, err := user.ToPublic(
		service.Logger,
		service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		return nil, "", common.ErrURLSigning
	}

	// Generate JWT accessToken
	accessToken, err := service.generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return publicUser, accessToken, nil
}

func (service *AuthService) Login(body *LoginBody) (*models.PublicUser, string, error) {
	user, err := service.UserService.GetUserByEmail(body.Email)
	if err != nil {
		return nil, "", common.ErrInvalidCredentials
	}

	// Compare password with hashed
	if !common.CheckHashedPassword(body.Password, user.Password) {
		return nil, "", common.ErrInvalidCredentials
	}

	publicUser, err := user.ToPublic(
		service.Logger,
		service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		return nil, "", common.ErrURLSigning
	}

	// Generate JWT accessToken
	accessToken, err := service.generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return publicUser, accessToken, nil
}

func (service *AuthService) GetResetPwdMail(email string) error {
	// Check if email actually existed
	user, err := service.UserService.GetUserByEmail(email)
	if err != nil {
		return err
	}

	resetPwdToken, err := service.generateActionToken(user.Email, time.Minute*10)
	if err != nil {
		return err
	}

	err = service.MailService.SendResetPwd(user, resetPwdToken)
	return err
}

func (service *AuthService) ResetPwd(body *ResetPwdBody) error {
	// Parse resetPwdToken
	resetPwdToken, err := common.ParseJWTToken(body.ResetPwdToken, service.AppConfig.JWTSecret)
	if err != nil {
		service.Logger.Debug("Reset password token parse failed", zap.Error(err))
		if errors.Is(err, jwt.ErrTokenExpired) {
			return common.ErrActionTokenExpired
		}
		return common.ErrActionTokenParsing
	}
	if !resetPwdToken.Valid {
		return common.ErrInvalidActionToken
	}

	// Get email from accessToken claims
	claims, ok := resetPwdToken.Claims.(jwt.MapClaims)
	if !ok {
		service.Logger.Debug("Reset password token type assertion failed")
		return common.ErrActionTokenClaimsRetrieval
	}
	email, ok := claims["email"].(string)
	if email == "" || !ok {
		service.Logger.Debug("Email invalid or missing", zap.String("email", email))
		return common.ErrActionTokenClaimsRetrieval
	}

	// Hash and updaate password
	err = service.UserService.UpdateUserPwdByEmail(email, body.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

// ======================== HELPER METHODS ========================

func (service *AuthService) generateActionToken(
	email string,
	expiresIn time.Duration,
) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}

	signedToken, err := common.GenerateJTWToken(claims,
		service.AppConfig.JWTSecret,
		expiresIn)
	if err != nil {
		service.Logger.Error("JWT action token generation failed", zap.Error(err))
		return "", common.ErrTokenGeneration
	}

	return signedToken, nil
}

func (service *AuthService) generateAccessToken(
	userID uuid.UUID,
	role types.UserRole,
) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}

	signedToken, err := common.GenerateJTWToken(claims,
		service.AppConfig.JWTSecret,
		time.Duration(service.AppConfig.JWTExpiresIn)*time.Hour)
	if err != nil {
		service.Logger.Error("JWT access token generation failed", zap.Error(err))
		return "", common.ErrTokenGeneration
	}

	return signedToken, nil
}
