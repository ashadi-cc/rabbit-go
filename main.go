package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	defer recoverAction()
	//setup rabbitmq
	conn, err := amqp.Dial("amqp://user:user@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ %s", err.Error())
	}
	defer conn.Close()

	action := enterChoice()
	if action == "1" {
		sendMessage(conn)
	}

	if action == "2" {
		receiveMessage(conn)
	}
}

func enterChoice() string {
	fmt.Println("Select RabbitMQ service:")
	fmt.Println("1. Producer")
	fmt.Println("2. Consumer")

	var action string
	fmt.Print("Enter action: ")
	fmt.Scanln(&action)
	switch action {
	case "1":
		fmt.Println("Run Producer")
	case "2":
		fmt.Println("Run Consumer")
	default:
		panic(fmt.Sprintf("Action not recognize"))
	}

	return action
}

func recoverAction() {
	action := recover()
	str, ok := action.(string)
	if ok && str == "Action not recognize" {
		fmt.Println(str)
		enterChoice()
	}
}
