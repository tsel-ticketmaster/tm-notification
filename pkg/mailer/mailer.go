package mailer

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Content Type
const (
	ContentTypePlaintext = "text/plain"
	ContentTypeHTML      = "text/html"
)

// Mailer error
var (
	ErrNoMessage   = fmt.Errorf("Mailer: No message to be sent")
	ErrNoRecipient = fmt.Errorf("Mailer: No email recipient")
)

// Recepient is data type for the mail recipient.
type Recepient struct {
	Address string
	Name    string
}

// MessageBody is the message body and the content type.
type MessageBody struct {
	ContentType string
	Body        []byte
}

// Message is a message to be sent to the mail server.
type Message struct {
	From        string
	To          []Recepient
	CC          []Recepient
	Subject     string
	MessageBody MessageBody
}

// Mailer is collection of behavior of mailer.
type Mailer interface {
	Send(ctx context.Context, messages ...Message) (err error)
}

type unimplementMailer struct {
	logger        *logrus.Logger
	defaultSender string
}

func (um *unimplementMailer) Send(ctx context.Context, messages ...Message) (err error) {
	for i, message := range messages {
		for j, recipient := range message.To {
			from := um.defaultSender

			if from != "" {
				from = recipient.Address
			}

			um.logger.WithContext(ctx).WithFields(logrus.Fields{
				"no":            fmt.Sprintf("%d.%d", i, j),
				"email.subject": message.Subject,
				"email.from":    from,
				"email.to":      recipient.Address,
			}).Info("fake sending email")
		}
	}

	return nil
}
