package customer

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-notification/pkg/errors"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailtemplate"
	"github.com/tsel-ticketmaster/tm-notification/pkg/status"
)

type CustomerUseCase interface {
	OnSignUp(ctx context.Context, event SignUpEvent) error
	OnChangeEmail(ctx context.Context, event ChangeEmailEvent) error
}

type CustomerUseCaseProperty struct {
	AppName     string
	Logger      *logrus.Logger
	EmailSender string
	Mailer      mailer.Mailer
}

type customerUseCase struct {
	appName     string
	logger      *logrus.Logger
	emailSender string
	mailer      mailer.Mailer
}

func NewCustomerUseCase(props CustomerUseCaseProperty) CustomerUseCase {
	return &customerUseCase{
		appName:     props.AppName,
		logger:      props.Logger,
		emailSender: props.EmailSender,
		mailer:      props.Mailer,
	}
}

// OnChangeEmail implements CustomerUseCase.
func (u *customerUseCase) OnChangeEmail(ctx context.Context, event ChangeEmailEvent) error {
	return nil
}

// OnSignUp implements CustomerUseCase.
func (u *customerUseCase) OnSignUp(ctx context.Context, event SignUpEvent) error {
	emailSubject := "Customer Verification"
	recipients := make([]mailer.Recepient, 1)
	recipients[0] = mailer.Recepient{
		Address: event.Email,
		Name:    event.Name,
	}

	data := &mailtemplate.VerificationEmailData{
		RecipientName:    event.Name,
		VerificationLink: event.VerificationLink,
	}

	mt := mailtemplate.NewCustomerVerificationTemplate()
	mtBuff := mt.Populate(data)

	if err := u.mailer.Send(context.TODO(), mailer.Message{
		From:    u.emailSender,
		To:      recipients,
		Subject: emailSubject,
		MessageBody: mailer.MessageBody{
			ContentType: "text/html",
			Body:        mtBuff.Bytes(),
		},
	}); err != nil {
		u.logger.WithContext(ctx).WithError(err).WithField("event", event).Error()
		return errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, err.Error())
	}

	return nil
}
