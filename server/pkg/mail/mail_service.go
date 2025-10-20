package mailfx

import (
	"fmt"
	"html/template"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
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
}

const (
	sender  = "noreply@touchgrassscheduler.com"
	appName = "Touch-Grass-Scheduler"
)

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

	return &MailService{
		AppConfig:                   params.AppConfig,
		Logger:                      params.Logger,
		MailClient:                  params.MailClient,
		RegistrationWarningTpl:      registrationWarningTpl,
		RegistrationVerificationTpl: registrationVerificationTpl,
	}
}

func (service *MailService) SendRegistrationWarning(user *models.User) error {
	msg := gomail.NewMsg()
	msg.To(user.Email)
	msg.From(sender)
	msg.Subject("Did you try to sign up for Touch-Grass-Scheduler?")

	// NOTE: The client endpoint is hard-coded
	data := &struct {
		UserFirstName     string
		UserEmail         string
		AppName           string
		ForgotPasswordURL string
	}{
		UserFirstName: user.FirstName,
		UserEmail:     user.Email,
		AppName:       appName,
		ForgotPasswordURL: fmt.Sprintf("%s:%d/forgot-password",
			service.AppConfig.ClientURL,
			service.AppConfig.ClientPort,
		),
	}

	if err := msg.SetBodyHTMLTemplate(service.RegistrationWarningTpl, data); err != nil {
		service.Logger.Error("Error setting HTML body to the message", zap.Error(err))
		return common.ErrMailHTMLSetting
	}

	if err := service.MailClient.DialAndSend(msg); err != nil {
		service.Logger.Error("Error sending message", zap.Error(err))
		return common.ErrMailSending
	}

	return nil
}

func (service *MailService) SendRegistrationVerification(email string, registrationToken string) error {
	msg := gomail.NewMsg()
	msg.To(email)
	msg.From(sender)
	msg.Subject(fmt.Sprintf("Complete your %s registration", appName))

	// NOTE: The client endpoint is hard-coded
	data := &struct {
		AppName         string
		JWTExpiresIn    int
		RegistrationURL string
	}{
		AppName:      appName,
		JWTExpiresIn: service.AppConfig.JWTExpiresIn,
		RegistrationURL: fmt.Sprintf("%s:%d/register/%s",
			service.AppConfig.ClientURL,
			service.AppConfig.ClientPort,
			registrationToken,
		),
	}

	if err := msg.SetBodyHTMLTemplate(service.RegistrationVerificationTpl, data); err != nil {
		service.Logger.Error("Error setting HTML body to the message", zap.Error(err))
		return common.ErrMailHTMLSetting
	}

	if err := service.MailClient.DialAndSend(msg); err != nil {
		service.Logger.Error("Error sending message", zap.Error(err))
		return common.ErrMailSending
	}

	return nil
}
