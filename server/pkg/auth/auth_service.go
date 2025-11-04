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

	if user != nil {
		// Send warning email if user already existed
		err = service.MailService.SendRegistrationWarning(user)
		if err != nil {
			return err
		}
	} else {
		// Send verification email if it is new user
		var registrationToken string
		registrationToken, err = service.generateActionToken(email,
			time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
		if err != nil {
			return err
		}

		// WARN: This is only for testing
		service.Logger.Debug("Temporary", zap.String("registrationToken", registrationToken))
		return nil

		err = service.MailService.SendRegistrationVerification(email, registrationToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *AuthService) Register(registrationTokenString string, body *RegisterBody) (*models.PublicUser, string, error) {
	// TODO: Add logic to check if SchoolNumber is valid

	// Parse registerToken
	registrationToken, err := common.ParseJWTToken(registrationTokenString, service.AppConfig.JWTSecret)
	if err != nil {
		service.Logger.Debug("Error parsing registrationToken", zap.Error(err))
		return nil, "", common.ErrActionTokenParsing
	}

	// Get email from accessToken claims
	claims, ok := registrationToken.Claims.(jwt.MapClaims)
	if !ok || !registrationToken.Valid {
		service.Logger.Debug("Error getting claims from registrationToken")
		return nil, "", common.ErrActionTokenClaimsGetting
	}
	email, ok := claims["email"].(string)
	if email == "" || !ok {
		service.Logger.Debug("Error getting email from claims of registrationToken")
		return nil, "", common.ErrActionTokenClaimsGetting
	}

	// Check if email is valid
	if _, err = mail.ParseAddress(email); err != nil {
		service.Logger.Debug("Error email is invalid")
		return nil, "", common.ErrVariableParsing
	}

	// Create new user
	user := body.ToUserModel()
	user.Email = email
	if err := service.UserService.CreateUser(user); err != nil {
		return nil, "", err
	}
	publicUser, err := user.ToPublic(service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		service.Logger.Error("Error getting signed URL for user's avatar", zap.Error(err))
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

	publicUser, err := user.ToPublic(service.StorageClient,
		service.AppConfig.StorageBucketName,
		time.Hour*time.Duration(service.AppConfig.JWTExpiresIn))
	if err != nil {
		service.Logger.Error("Error getting signed URL for user's avatar", zap.Error(err))
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
		if errors.Is(err, common.ErrUserNotFound) {
			service.Logger.Debug("")
			return err
		}

		service.Logger.Error("", zap.Error(err))
		return common.ErrDatabase
	}

	resetPwdToken, err := service.generateActionToken(user.Email, time.Minute*5)
	if err != nil {
		return err
	}

	err = service.MailService.SendResetPwd(user, resetPwdToken)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) ResetPwd(body *ResetPwdBody) error {
	// Parse resetPwdToken
	resetPwdToken, err := common.ParseJWTToken(body.ResetPwdToken, service.AppConfig.JWTSecret)
	if err != nil {
		service.Logger.Debug("Error parsing resetPwdToken", zap.Error(err))
		return common.ErrActionTokenParsing
	}

	// Get email from accessToken claims
	claims, ok := resetPwdToken.Claims.(jwt.MapClaims)
	if !ok || !resetPwdToken.Valid {
		service.Logger.Debug("Error getting claims from resetPwdToken")
		return common.ErrActionTokenClaimsGetting
	}
	email, ok := claims["email"].(string)
	if email == "" || !ok {
		service.Logger.Debug("Error getting email from claims of resetPwdToken")
		return common.ErrActionTokenClaimsGetting
	}

	// Hash and updaate password
	err = service.UserService.UpdateUserPwdByEmail(email, body.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

// ======================== HELPER METHODS ========================

func (service *AuthService) generateActionToken(email string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}

	signedToken, err := common.GenerateJTWToken(claims,
		service.AppConfig.JWTSecret,
		expiresIn)
	if err != nil {
		service.Logger.Error("Internal error while signing JWT actionToken:", zap.Error(err))
		return "", common.ErrTokenGeneration
	}

	return signedToken, nil
}

func (service *AuthService) generateAccessToken(userID uuid.UUID, role types.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}

	signedToken, err := common.GenerateJTWToken(claims,
		service.AppConfig.JWTSecret,
		time.Duration(service.AppConfig.JWTExpiresIn)*time.Hour)
	if err != nil {
		service.Logger.Error("Internal error while signing JWT accessToken:", zap.Error(err))
		return "", common.ErrTokenGeneration
	}

	return signedToken, nil
}
