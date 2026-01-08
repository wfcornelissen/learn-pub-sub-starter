package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("Failed to establish connection:/n%v/n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection successfully established!")

	connChannel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to create a channel for the connection:/n%v/n", err)
		return
	}

	err = pubsub.PublishJSON(connChannel, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("Failed to publish:/n%v/n", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
