package ticket

import "time"

type AcquireTicketEvent struct {
	ID                   int64
	Number               string
	EventID              string
	ShowID               string
	Tier                 string
	TicketStockID        string
	EventName            string
	ShowVenue            string
	ShowType             string
	ShowCountry          string
	ShowCity             string
	ShowFormattedAddress string
	ShowTime             time.Time
	CustomerName         string
	CustomerEmail        string
	CustomerID           int64
	CreatedAt            time.Time
	OrderID              string
}
