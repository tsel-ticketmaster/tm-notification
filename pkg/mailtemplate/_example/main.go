package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailtemplate/v2"
	"gopkg.in/gomail.v2"
)

func main() {
	data := &mailtemplate.ExampleData{
		RecipientName: "JINDAN DUKUN KEBAL",
	}
	mt := mailtemplate.NewExampleTemplate()
	mtBuff := mt.Populate(data)

	gomailDialer := gomail.NewDialer(
		"smtp-relay.sendinblue.com", 587,
		"support@my-pertamina.id", "xsmtpsib-3e9a775d05204c19b866c2709cecddde052e0c067aa5861b0e42709961d5b6ef-mjzTdvaJHGrqtRED",
	)

	m := mailer.NewGomailAdapter(logrus.New(), `"KANG UDIN" <no-reply@my-pertamina.id>`, gomailDialer)
	err := m.Send(context.TODO(), mailer.Message{
		From: `"KANG UDIN" <no-reply@my-pertamina.id>`,
		To: []mailer.Recepient{
			{
				Address: "ijal.alfarizi@gmail.com",
				Name:    "Jindan",
			},
		},
		Subject: "UNDANGAN MEETING PERDUKUNAN",
		MessageBody: mailer.MessageBody{
			ContentType: "text/html",
			Body:        mtBuff.Bytes(),
		},
	})

	if err != nil {
		panic(err)
	}
}
