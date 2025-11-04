package mailfx

import (
	"fmt"
	"html/template"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/endpoints"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	gomail "github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MailServiceParams struct {
	fx.In
	AppConfig  *configfx.AppConfig
	Logger     *zap.Logger
	MailClient *gomail.Client
}

type MailService struct {
	AppConfig                   *configfx.AppConfig
	Logger                      *zap.Logger
	MailClient                  *gomail.Client
	RegistrationWarningTpl      *template.Template
	RegistrationVerificationTpl *template.Template
	ResetPwdTpl                 *template.Template
}

const (
	sender  = "noreply@touchgrassscheduler.com"
	appName = "Touch-Grass-Scheduler"
)

// ======================== METHODS ========================

func NewMailService(params MailServiceParams) *MailService {
	registrationWarningTpl, err := template.
		ParseFiles("pkg/mail/templates/registration_warning.html")
	if err != nil {
		params.Logger.Fatal("Error parsing Registration Warning Template", zap.Error(err))
	}

	registrationVerificationTpl, err := template.
		ParseFiles("pkg/mail/templates/registration_verification.html")
	if err != nil {
		params.Logger.Fatal("Error parsing Registration Verification Template", zap.Error(err))
	}

	resetPwdTpl, err := template.ParseFiles("pkg/mail/templates/reset_password.html")
	if err != nil {
		params.Logger.Fatal("Error parsing Reset Password Template", zap.Error(err))
	}

	return &MailService{
		AppConfig:                   params.AppConfig,
		Logger:                      params.Logger,
		MailClient:                  params.MailClient,
		RegistrationWarningTpl:      registrationWarningTpl,
		RegistrationVerificationTpl: registrationVerificationTpl,
		ResetPwdTpl:                 resetPwdTpl,
	}
}

func (service *MailService) SendRegistrationWarning(user *models.User) error {
	subject := "Did you try to sign up for Touch-Grass-Scheduler?"

	data := &struct {
		UserFirstName     string
		UserEmail         string
		AppName           string
		ForgotPasswordURL string
	}{
		UserFirstName: user.FirstName,
		UserEmail:     user.Email,
		AppName:       appName,
		ForgotPasswordURL: fmt.Sprintf("%s/%s",
			service.AppConfig.ClientURL,
			endpoints.ClientForgotPwd,
		),
	}

	err := service.setBodyAndSend(user.Email, sender, subject, service.RegistrationWarningTpl, data)
	if err != nil {
		return err
	}

	return nil
}

func (service *MailService) SendRegistrationVerification(email string, registrationToken string) error {
	subject := fmt.Sprintf("Complete your %s registration", appName)

	data := &struct {
		AppName         string
		JWTExpiresIn    int
		RegistrationURL string
	}{
		AppName:      appName,
		JWTExpiresIn: service.AppConfig.JWTExpiresIn,
		RegistrationURL: fmt.Sprintf("%s/%s/%s",
			service.AppConfig.ClientURL,
			endpoints.ClientRegistrationVerification,
			registrationToken,
		),
	}

	err := service.setBodyAndSend(email, sender, subject, service.RegistrationVerificationTpl, data)
	if err != nil {
		return err
	}

	return nil
}

func (service *MailService) SendResetPwd(user *models.User, resetPwdToken string) error {
	subject := fmt.Sprintf("Reset your password on %s", appName)

	// NOTE: ExpireIn is hardcoded. Please set to sync with auth_service
	data := &struct {
		UserFirstName string
		AppName       string
		ExpiresIn     int
		ResetPwdURL   string
	}{
		UserFirstName: user.FirstName,
		AppName:       appName,
		ExpiresIn:     5, // minutes
		ResetPwdURL: fmt.Sprintf("%s/%s/%s",
			service.AppConfig.ClientURL,
			endpoints.ClientResetPwd,
			resetPwdToken,
		),
	}

	err := service.setBodyAndSend(user.Email, sender, subject, service.ResetPwdTpl, data)
	if err != nil {
		return err
	}

	return nil
}

// ======================== HELPER METHODS ========================

func (service *MailService) setBodyAndSend(reciever, sender, subject string, tpl *template.Template, data any) error {
	var err error

	msg := gomail.NewMsg()

	err = msg.To(reciever)
	if err != nil {
		service.Logger.Error("", zap.Error(err))
		return err
	}

	err = msg.From(sender)
	if err != nil {
		service.Logger.Error("", zap.Error(err))
		return err
	}

	msg.Subject(subject)

	if err := msg.SetBodyHTMLTemplate(tpl, data); err != nil {
		service.Logger.Error("Error setting HTML body to the message", zap.Error(err))
		return common.ErrMailHTMLSetting
	}

	if err := service.MailClient.DialAndSend(msg); err != nil {
		service.Logger.Error("Error sending message", zap.Error(err))
		return common.ErrMailSending
	}

	return nil
}
