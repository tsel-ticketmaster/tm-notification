package ticket

import (
	"context"
	"encoding/json"
	"fmt"

	ck "github.com/confluentinc/confluent-kafka-go/kafka"
)

type AcquireTicketEventHandler struct {
	TicketUseCase TicketUseCase
}

func (handler AcquireTicketEventHandler) Handle(ctx context.Context, msg interface{}) error {
	kafkaMessage, ok := msg.(*ck.Message)
	if !ok {
		return fmt.Errorf("invalid message provider")
	}

	event := AcquireTicketEvent{}
	json.Unmarshal(kafkaMessage.Value, &event)

	return handler.TicketUseCase.OnAcquireTicket(ctx, event)
}
