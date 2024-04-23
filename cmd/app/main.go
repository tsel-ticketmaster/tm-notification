package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tsel-ticketmaster/tm-notification/config"
	customerapp_customer "github.com/tsel-ticketmaster/tm-notification/internal/module/customerapp/customer"
	customerapp_ticket "github.com/tsel-ticketmaster/tm-notification/internal/module/customerapp/ticket"
	"github.com/tsel-ticketmaster/tm-notification/pkg/applogger"
	"github.com/tsel-ticketmaster/tm-notification/pkg/kafka"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
	"github.com/tsel-ticketmaster/tm-notification/pkg/middleware"
	"github.com/tsel-ticketmaster/tm-notification/pkg/monitoring"
	"github.com/tsel-ticketmaster/tm-notification/pkg/pubsub"
	"github.com/tsel-ticketmaster/tm-notification/pkg/response"
	"github.com/tsel-ticketmaster/tm-notification/pkg/server"
	"github.com/tsel-ticketmaster/tm-notification/pkg/status"
	"github.com/tsel-ticketmaster/tm-notification/pkg/validator"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

var (
	c           *config.Config
	CustomerApp string
	AdminApp    string
)

func init() {
	c = config.Get()
	AdminApp = fmt.Sprintf("%s/%s", c.Application.Name, "adminapp")
	CustomerApp = fmt.Sprintf("%s/%s", c.Application.Name, "customerapp")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := applogger.GetLogrus()

	mon := monitoring.NewOpenTelemetry(
		c.Application.Name,
		c.Application.Environment,
		c.GCP.ProjectID,
	)

	mon.Start(ctx)

	cloudstorage, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(c.GCP.ServiceAccount))
	if err != nil {
		logger.WithError(err).Error()
	}

	_ = validator.Get()

	gomailDialer := gomail.NewDialer(
		c.Mailer.SMTP.Host, c.Mailer.SMTP.Port,
		c.Mailer.SMTP.Username, c.Mailer.SMTP.Password,
	)
	gomailAdapter := mailer.NewGomailAdapter(logger, c.Mailer.Sender, gomailDialer, true)

	router := mux.NewRouter()
	router.Use(
		otelmux.Middleware(c.Application.Name),
		middleware.HTTPResponseTraceInjection,
		middleware.NewHTTPRequestLogger(logger, c.Application.Debug).Middleware,
	)
	router.HandleFunc("/tm-notification", healthCheck).Methods(http.MethodGet)

	// admin's app

	// customer's app
	customerappCustomerUseCase := customerapp_customer.NewCustomerUseCase(customerapp_customer.CustomerUseCaseProperty{
		AppName:     CustomerApp,
		Logger:      logger,
		EmailSender: c.Mailer.Sender,
		Mailer:      gomailAdapter,
	})
	customerSignUpSubscriber := pubsub.SubscriberFromConfluentKafkaConsumer(pubsub.ConfluentKafkaConsumerProperty{
		Logger: logger,
		Topic:  "customer-sign-up",
		EventHandler: customerapp_customer.SignUpEventHandler{
			CustomerUseCase: customerappCustomerUseCase,
		},
		Consumer: kafka.NewConsumer(CustomerApp, false),
	})
	customerSignUpSubscriber.Subscribe()

	customerappTicketUseCase := customerapp_ticket.NewTicketUseCase(customerapp_ticket.TicketUseCaseProperty{
		AppName:      CustomerApp,
		Logger:       logger,
		EmailSender:  c.Mailer.Sender,
		Mailer:       gomailAdapter,
		CloudStorage: cloudstorage,
	})
	customerappAqcuireTicketSubscriber := pubsub.SubscriberFromConfluentKafkaConsumer(pubsub.ConfluentKafkaConsumerProperty{
		Logger: logger,
		Topic:  "acquire-ticket",
		EventHandler: customerapp_ticket.AcquireTicketEventHandler{
			TicketUseCase: customerappTicketUseCase,
		},
		Consumer: kafka.NewConsumer(CustomerApp, false),
	})
	customerappAqcuireTicketSubscriber.Subscribe()

	handler := middleware.SetChain(
		router,
		cors.New(cors.Options{
			AllowedOrigins:   c.CORS.AllowedOrigins,
			AllowedMethods:   c.CORS.AllowedMethods,
			AllowedHeaders:   c.CORS.AllowedHeaders,
			ExposedHeaders:   c.CORS.ExposedHeaders,
			MaxAge:           c.CORS.MaxAge,
			AllowCredentials: c.CORS.AllowCredentials,
		}).Handler,
	)

	srv := &server.Server{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", c.Application.Port),
			Handler: handler,
		},
		Logger: logger,
	}

	go func() {
		srv.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	srv.Shutdown(ctx)
	customerappAqcuireTicketSubscriber.Close()
	customerSignUpSubscriber.Close()
	mon.Stop(ctx)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, response.RESTEnvelope{
		Status:  status.OK,
		Message: fmt.Sprintf("%s is running properly", os.Getenv("APP_NAME")),
		Data:    nil,
		Meta:    nil,
	})
}
