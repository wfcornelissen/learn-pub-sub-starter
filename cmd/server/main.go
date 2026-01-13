package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("Failed to establish connection:\n%v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection successfully established!")

	connChannel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to create a channel for the connection:\n%v\n", err)
		return
	}

	err = pubsub.PublishJSON(connChannel, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("Failed to publish:\n%v\n", err)
		return
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "game_logs", "game_logs.*", pubsub.Durable)
	if err != nil {
		fmt.Printf("Failed to declare and bind queue:\n%v\n", err)
		return
	}

	gamelogic.PrintServerHelp()

loop:
	for true {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		}
		switch strings.ToLower(cmd[0]) {
		case "pause":
			fmt.Println("Sending pause")
			err = pubsub.PublishJSON(connChannel, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Printf("Failed to publish:/n%v/n", err)
				return
			}
		case "resume":
			fmt.Println("Sending resume")
			err = pubsub.PublishJSON(connChannel, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Printf("Failed to publish:/n%v/n", err)
				return
			}
		case "quit":
			fmt.Println("Exiting")
			break loop
		default:
			fmt.Println("Unknown command")
			continue loop
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
