package main

import (
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Define RabbitMQ server URL
	amqpSeverURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection
	connectRabbitMQ, err := amqp.Dial(amqpSeverURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

}
