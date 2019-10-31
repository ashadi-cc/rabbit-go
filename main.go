package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	//setup rabbitmq
	conn, err := amqp.Dial("amqp://user:user@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ %s", err.Error())
	}
	defer conn.Close()
	defer recoverAction(conn)
	enterChoice(conn)

}

func enterChoice(conn *amqp.Connection) {
	fmt.Println("Select RabbitMQ service:")
	fmt.Println("1. Producer")
	fmt.Println("2. Consumer")
	fmt.Println("3. Direct Producer")
	fmt.Println("4. Direct Consumer")
	fmt.Println("5. fanout Producer")
	fmt.Println("6. fanout Consumer")

	var action string
	fmt.Print("Enter action: ")
	fmt.Scanln(&action)
	switch action {
	case "1":
		fmt.Println("Run Producer")
		sendMessage(conn)
	case "2":
		fmt.Println("Run Consumer")
		receiveMessage(conn)
	case "3":
		fmt.Println("Run direct producer")
		sendMessageDirect(conn)
	case "4":
		fmt.Println("run direct consumer")
		receiveMessageDirect(conn)
	case "5":
		fmt.Println("Run fanout producer")
		sendMessageFannout(conn)
	case "6":
		fmt.Println("run fanout consumer")
		receiveMessageFannout(conn)
	default:
		panic(fmt.Sprintf("Action not recognize"))
	}
}

func recoverAction(conn *amqp.Connection) {
	action := recover()
	str, ok := action.(string)
	if ok && str == "Action not recognize" {
		fmt.Println(str)
		enterChoice(conn)
	}
}
