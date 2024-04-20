package mailtemplate

import "bytes"

// Data is an abstraction of mail data.
type Data interface {
	Get() interface{}
}

// MailTemplate is an abstraction of mail template.
type MailTemplate interface {
	Populate(data Data) (buff *bytes.Buffer)
}
