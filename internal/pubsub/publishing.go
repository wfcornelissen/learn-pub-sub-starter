package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType struct {
	Durable   bool
	Transient bool
}

func CreateQueueType(qtpe string) (SimpleQueueType, error) {
	if strings.ToLower(qtpe) == "durable" {
		return SimpleQueueType{Durable: true, Transient: false}, nil
	} else if strings.ToLower(qtpe) == "transient" {
		return SimpleQueueType{Durable: false, Transient: true}, nil
	}
	return SimpleQueueType{}, fmt.Errorf("Invalid type passed")
}

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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	fmt.Println("Starting Peril client...")
	channel, err := conn.Channel()
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("Error creating channel:\n%v\n", err)
	}

	queue, err := channel.QueueDeclare(queueName, queueType.Durable, queueType.Transient, queueType.Transient, false, nil)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("Error creating queue:\n%v\n", err)

	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("Error binding queue:\n%v\n", err)
	}

	return channel, queue, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("Error verifying queue existence/bind:\n%v\n", err)
	}
	deliveries, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving delivery structs:\n%v\n", err)
	}

	go func() {
		for delivery := range deliveries {
			var msg T
			err := json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				fmt.Printf("Error retrieving delivery structs:\n%v\n", err)
			}
			handler(msg)
			err = delivery.Ack(false)
			if err != nil {
				fmt.Printf("Failed to acknowledge delivery:\n%v\n", err)
			}
		}
	}()
	return nil

}
