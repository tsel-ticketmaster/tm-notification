package mailtemplate

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed html/acquired_ticket_notification.html
var acquiredTicketNotificationTemplate embed.FS

// AcquiredTicketNotificationTemplate is a concrete struct of MailTemplate
type AcquiredTicketNotificationTemplate struct {
	rawBuff []byte
	raw     string
}

// NewCustomerVerificationTemplate is a constructor.
func NewAcquiredTicketNotificationTemplate() MailTemplate {
	rawBuff, _ := acquiredTicketNotificationTemplate.ReadFile("html/acquired_ticket_notification.html")
	return &AcquiredTicketNotificationTemplate{
		rawBuff: rawBuff,
		raw:     string(rawBuff),
	}
}

// Populate will populate the template with data.
//
// The `Data` must contain fields:
//
// `RecipientName` As the name of recipient.
//
// For example:
//
//	// set the arguments
//	data := &ExampleData{}
//	data.RecipientName = "John Doe"
//
//	// instantiate the template
//	template := NewExampleTemplate()
//	tbuff := template.Populate(data)
func (et *AcquiredTicketNotificationTemplate) Populate(data Data) (buff *bytes.Buffer) {
	buff = new(bytes.Buffer)
	t, _ := template.New("acquired-ticket-notification").Parse(et.raw)
	t.Execute(buff, data)
	return
}
