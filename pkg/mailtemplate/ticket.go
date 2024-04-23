package mailtemplate

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed html/ticket.html
var ticketTemplate embed.FS

// TicketTemplate is a concrete struct of MailTemplate
type TicketTemplate struct {
	rawBuff []byte
	raw     string
}

// NewTicketTemplate is a constructor.
func NewTicketTemplate() MailTemplate {
	rawBuff, _ := ticketTemplate.ReadFile("html/ticket.html")
	return &TicketTemplate{
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
func (et *TicketTemplate) Populate(data Data) (buff *bytes.Buffer) {
	buff = new(bytes.Buffer)
	t, _ := template.New("acquired-ticket-notification").Parse(et.raw)
	t.Execute(buff, data)
	return
}
