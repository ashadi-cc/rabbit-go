package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func sendMessage(connection *amqp.Connection) {
	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	defer channel.Close()

	// We create an exahange that will bind to the queue to send and receive messages
	// see this for kind of type https://medium.com/faun/different-types-of-rabbitmq-exchanges-9fefd740505d
	err = channel.ExchangeDeclare("events2", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	loop, index := true, 1
	for loop {

		select {
		case <-ctx.Done():
			loop = false
		case <-time.After(time.Second):
			// We create a message to be sent to the queue.
			// It has to be an instance of the aqmp publishing struct
			log.Println("send Message ", fmt.Sprintf("Hello World %d", index))
			message := amqp.Publishing{
				Body: []byte(fmt.Sprintf("Hello World %d", index)),
			}

			// We publish the message to the exahange we created earlier
			err = channel.Publish("events2", "random-key", false, false, message)

			if err != nil {
				panic("error publishing a message to the queue:" + err.Error())
			}
		}
		index++
	}

}

func receiveMessage(connection *amqp.Connection) {
	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	defer channel.Close()

	// We create an exahange that will bind to the queue to send and receive messages
	// see this for kind of type https://medium.com/faun/different-types-of-rabbitmq-exchanges-9fefd740505d
	err = channel.ExchangeDeclare("events2", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}
	// We create a queue named Test
	_, err = channel.QueueDeclare("test", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	//with any routing key #
	err = channel.QueueBind("test", "#", "events2", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	// We consume data in the queue named test using the channel we created in go.
	msgs, err := channel.Consume("test", "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	// We loop through the messages in the queue and print them to the console.
	// The msgs will be a go channel, not an amqp channel
	for msg := range msgs {
		//print the message to the console
		fmt.Println("message received: " + string(msg.Body))
		// Acknowledge that we have received the message so it can be removed from the queue
		msg.Ack(false)
	}

}
