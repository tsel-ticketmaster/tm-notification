package mailtemplate

type VerificationEmailData struct {
	RecipientName    string
	VerificationLink string
}

func (d VerificationEmailData) Get() interface{} {
	return d
}

type AcquiredTicketNotificationData struct {
	CustomerName  string
	TicketPDFLink string
}

func (d AcquiredTicketNotificationData) Get() interface{} {
	return d
}

type TicketData struct {
	CustomerName string
	EventName    string
	Venue        string
	Country      string
	City         string
	Tier         string
	TicketNumber string
	DateTime     string
}

func (d TicketData) Get() interface{} {
	return d
}
