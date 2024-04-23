package ticket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-notification/pkg/errors"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailtemplate"
	"github.com/tsel-ticketmaster/tm-notification/pkg/status"
)

type TicketUseCase interface {
	OnAcquireTicket(ctx context.Context, e AcquireTicketEvent) error
}

type TicketUseCaseProperty struct {
	AppName      string
	Logger       *logrus.Logger
	EmailSender  string
	Mailer       mailer.Mailer
	CloudStorage *storage.Client
}

type ticketUseCase struct {
	appName      string
	logger       *logrus.Logger
	emailSender  string
	mailer       mailer.Mailer
	cloudstorage *storage.Client
}

// OnAcquireTicket implements TicketUseCase.
func (u *ticketUseCase) OnAcquireTicket(ctx context.Context, e AcquireTicketEvent) error {
	tmt := mailtemplate.NewTicketTemplate()
	tmtBuff := tmt.Populate(&mailtemplate.TicketData{
		CustomerName: e.CustomerName,
		EventName:    e.EventName,
		Venue:        e.ShowVenue,
		Country:      e.ShowCountry,
		City:         e.ShowCity,
		Tier:         e.Tier,
		TicketNumber: e.Number,
		DateTime:     e.ShowTime.Format(time.DateTime),
	})

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var pdfBytes []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				u.logger.WithContext(ctx).WithError(err).Error()
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, tmtBuff.String()).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			printResult, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				u.logger.WithContext(ctx).WithError(err).Error()
				return err
			}

			pdfBytes = printResult
			return nil
		}),
	)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error()
		return err
	}

	filename := fmt.Sprintf("%s.pdf", e.Number)
	w := u.cloudstorage.Bucket("tsel-ticketmaster").Object(filename).NewWriter(ctx)
	w.ChunkSize = 0
	w.ContentType = "application/pdf"

	if _, err = io.Copy(w, bytes.NewBuffer(pdfBytes)); err != nil {
		u.logger.WithContext(ctx).WithError(err).Error()
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := w.Close(); err != nil {
		u.logger.WithContext(ctx).WithError(err).Error()
		return fmt.Errorf("Writer.Close: %w", err)
	}

	emailSubject := "Acquired Ticket"
	recipients := make([]mailer.Recepient, 1)
	recipients[0] = mailer.Recepient{
		Address: e.CustomerEmail,
		Name:    e.CustomerName,
	}

	pdfUrl := fmt.Sprintf("https://storage.googleapis.com/tsel-ticketmaster/%s", filename)
	data := &mailtemplate.AcquiredTicketNotificationData{
		CustomerName:  e.CustomerName,
		TicketPDFLink: pdfUrl,
	}

	mt := mailtemplate.NewAcquiredTicketNotificationTemplate()
	mtBuff := mt.Populate(data)

	if err := u.mailer.Send(context.TODO(), mailer.Message{
		From:    u.emailSender,
		To:      recipients,
		Subject: emailSubject,
		MessageBody: mailer.MessageBody{
			ContentType: "text/html",
			Body:        mtBuff.Bytes(),
		},
	}); err != nil {
		u.logger.WithContext(ctx).WithError(err).WithField("event", e).Error()
		return errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, err.Error())
	}

	return nil
}

func NewTicketUseCase(props TicketUseCaseProperty) TicketUseCase {
	return &ticketUseCase{
		appName:      props.AppName,
		logger:       props.Logger,
		emailSender:  props.EmailSender,
		mailer:       props.Mailer,
		cloudstorage: props.CloudStorage,
	}
}
