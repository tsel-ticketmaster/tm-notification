package customer

import (
	"context"
	"encoding/json"
	"fmt"

	ck "github.com/confluentinc/confluent-kafka-go/kafka"
)

type SignUpEventHandler struct {
	CustomerUseCase CustomerUseCase
}

func (handler SignUpEventHandler) Handle(ctx context.Context, msg interface{}) error {
	fmt.Println("masuk ee")
	kafkaMessage, ok := msg.(*ck.Message)
	if !ok {
		return fmt.Errorf("invalid message provider")
	}

	event := SignUpEvent{}
	json.Unmarshal(kafkaMessage.Value, &event)

	return handler.CustomerUseCase.OnSignUp(ctx, event)
}
