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
loop:
	for true {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		}
		switch strings.ToLower(cmd[0]) {
		case "spawn":
			//Validate length
			if len(cmd) < 3 {
				fmt.Println("Invalid command")
				continue loop
			}
			// Validate location
			locationBool := false
			for location := range gamelogic.GetAllLocations() {
				if cmd[1] == string(location) {
					locationBool = true
				}
			}
			// Validate unit
			unitBool := false
			for unit := range gamelogic.GetAllRanks() {
				if cmd[2] == string(unit) {
					unitBool = true
				}
			}
			// Continue loop if validation fails
			if !unitBool {
				fmt.Println("Invalid unit")
				continue loop
			}
			if !locationBool {
				fmt.Println("Invalid location")
				continue loop
			}

			err = ngs.CommandSpawn(cmd[1:])
			if err != nil {
				fmt.Printf("Error spawning unit:/\n%v\n", err)
			}
		case "move":
			//Validate length
			if len(cmd) < 3 {
				fmt.Println("Invalid command")
				continue loop
			}
			// Validate location
			locationBool := false
			for location := range gamelogic.GetAllLocations() {
				if cmd[1] == string(location) {
					locationBool = true
				}
			}
			// Validate unit
			unitBool := false
			for unit := range gamelogic.GetAllRanks() {
				if cmd[2] == string(unit) {
					unitBool = true
				}
			}
			// Continue loop if validation fails
			if !unitBool {
				fmt.Println("Invalid unit")
				continue loop
			}
			if !locationBool {
				fmt.Println("Invalid location")
				continue loop
			}
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
