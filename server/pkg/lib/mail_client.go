package libfx

import (
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	gomail "github.com/wneessen/go-mail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MailClientParams struct {
	fx.In
	AppParams *configfx.AppConfig
	Logger    *zap.Logger
}

func NewMailClient(params MailClientParams) *gomail.Client {
	client, err := gomail.NewClient(
		params.AppParams.MailHost,
		gomail.WithPort(params.AppParams.MailPort),
		gomail.WithUsername(params.AppParams.MailUser),
		gomail.WithPassword(params.AppParams.MailPassword),
		gomail.WithTLSPolicy(gomail.TLSMandatory),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
	)
	if err != nil {
		params.Logger.Fatal("Error creating mail client:", zap.Error(err))
	}

	params.Logger.Debug("Mail client initialized successfully.")

	return client
}
