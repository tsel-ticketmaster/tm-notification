package mailtemplate

type VerificationEmailData struct {
	RecipientName    string
	VerificationLink string
}

func (d VerificationEmailData) Get() interface{} {
	return d
}
