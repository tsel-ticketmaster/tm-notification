package mailer

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"gopkg.in/gomail.v2"
)

// GomailDialer is an abstraction of gomail.Dialer.
type GomailDialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// GomailAdapter is a concrete struct of gomail adapter.
type GomailAdapter struct {
	logger        *logrus.Logger
	defaultSender string
	dialer        GomailDialer
}

// NewGomailAdapter is a constructor.
func NewGomailAdapter(logger *logrus.Logger, defaultSender string, dialer GomailDialer, active bool) Mailer {
	if logger == nil {
		logger = logrus.New()
	}

	if !active {
		return &unimplementMailer{
			logger:        logger,
			defaultSender: defaultSender,
		}
	}

	return &GomailAdapter{
		logger:        logger,
		defaultSender: defaultSender,
		dialer:        dialer,
	}
}

// Send will send the email.
func (g *GomailAdapter) Send(ctx context.Context, messages ...Message) (err error) {
	tp := otel.GetTracerProvider()
	t := tp.Tracer("gomail")
	_, span := t.Start(ctx, "send")
	defer span.End()

	lengthOfMessages := len(messages)

	if lengthOfMessages < 1 {
		return ErrNoMessage
	}
	gomailMessages := make([]*gomail.Message, len(messages))

	for i, message := range messages {
		gm, err := g.composeGomailMessage(message)
		if err != nil {
			return err
		}

		gomailMessages[i] = gm
	}

	err = g.dialer.DialAndSend(gomailMessages...)
	return
}

func (g *GomailAdapter) setFrom(m Message, gm *gomail.Message) (err error) {
	sender := g.defaultSender
	if m.From != "" {
		sender = m.From
	}
	gm.SetHeader("From", sender)
	return
}

func (g *GomailAdapter) setRecipient(m Message, gm *gomail.Message) (err error) {
	// Set Recipient
	if len(m.To) > 1 { // what if recipients is more than one.
		addresses := make([]string, len(m.To))
		for i, recipient := range m.To {
			addresses[i] = recipient.Address
		}
		gm.SetHeader("To", addresses...)

	} else if len(m.To) == 1 { // what if it is single recipient.

		address := m.To[0]
		gm.SetAddressHeader("To", address.Address, address.Name)

	} else { // what if no recipient

		err = ErrNoRecipient
		return
	}
	return
}

func (g *GomailAdapter) setSubject(m Message, gm *gomail.Message) (err error) {
	gm.SetHeader("Subject", m.Subject)
	return
}

func (g *GomailAdapter) setCarbonCopy(m Message, gm *gomail.Message) (err error) {
	if len(m.CC) > 1 {
		CCs := make([]string, len(m.CC))
		for i, cc := range m.CC {
			CCs[i] = cc.Address
		}
		gm.SetHeader("Cc", CCs...)
	} else if len(m.CC) == 1 {
		cc := m.CC[0]
		gm.SetAddressHeader("Cc", cc.Address, cc.Name)
	}
	return
}

func (g *GomailAdapter) setBody(m Message, gm *gomail.Message) (err error) {
	gm.SetBody(m.MessageBody.ContentType, string(m.MessageBody.Body))
	return
}

func (g *GomailAdapter) composeGomailMessage(m Message) (gm *gomail.Message, err error) {
	gm = gomail.NewMessage()

	err = g.setRecipient(m, gm)
	if err != nil {
		gm = nil
		return
	}

	g.setFrom(m, gm)
	g.setSubject(m, gm)
	g.setCarbonCopy(m, gm)
	g.setBody(m, gm)

	return
}
