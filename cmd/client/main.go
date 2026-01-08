package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("Failed to establish connection:/n%v/n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection successfully established!")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("Error retrieving welcome message:/\n%v\n", err)
	}

	qtpe, err := pubsub.CreateQueueType("transient")
	if err != nil {
		fmt.Printf("Error creating queue type:/\n%v\n", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%v.%v", routing.PauseKey, userName),
		routing.PauseKey,
		qtpe)

	ngs := gamelogic.NewGameState(userName)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
