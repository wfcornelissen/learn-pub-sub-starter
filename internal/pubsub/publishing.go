package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	byteVal, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Failed to marshal value:/n%v/n", err)
	}

	newPub := amqp.Publishing{
		ContentType: "application/json",
		Body:        byteVal,
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, newPub)
	if err != nil {
		return fmt.Errorf("Failed to publish to channel:/n%v/n", err)
	}

	return nil
}
