package mailtemplate

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed html/customer_verification_template.html
var customerVerificationTemplate embed.FS

// CustomerVerificationTemplate is a concrete struct of MailTemplate
type CustomerVerificationTemplate struct {
	rawBuff []byte
	raw     string
}

// NewCustomerVerificationTemplate is a constructor.
func NewCustomerVerificationTemplate() MailTemplate {
	rawBuff, _ := customerVerificationTemplate.ReadFile("html/customer_verification_template.html")
	return &CustomerVerificationTemplate{
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
func (et *CustomerVerificationTemplate) Populate(data Data) (buff *bytes.Buffer) {
	buff = new(bytes.Buffer)
	t, _ := template.New("customer-verification").Parse(et.raw)
	t.Execute(buff, data)
	return
}
