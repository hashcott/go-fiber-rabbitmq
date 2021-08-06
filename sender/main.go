package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// Opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established

	channelRabbitMQ, err := connectRabbitMQ.Channel()

	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new fiber instance
	app := fiber.New()

	// Add middleware
	app.Use(
		logger.New(),
	)

	// Add route
	app.Get("/send", func(c *fiber.Ctx) error {
		// Create message to publish
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}
		// Attempt to publish a message to the queue.
		if err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		); err != nil {
			return err
		}
		return nil
	})
	log.Fatal(app.Listen(":3000"))
}
