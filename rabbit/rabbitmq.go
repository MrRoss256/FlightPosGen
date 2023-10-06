package rabbit

import (
	"context"
	"encoding/json"
	"fpg_types"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var conn *amqp.Connection

func Connect() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/guest")
	failOnError(err, "Failed to connect to RabbitMQ")
}

func Write(point *fpg_types.FlightLoc) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"flight", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(point)
	if err != nil {
		log.Error(err)
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Debugf("Sent %s @ %s\n", point.FlightDes, point.EventTime)
}
