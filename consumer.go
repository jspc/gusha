package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var (
	done chan error
)

type Incoming struct {
	Name    string
	BaseURL string
	Clients int
	Seconds int
	URLs    URLs
}

func Consumer() (err error) {
	conn, err := amqp.Dial(*amqpURI)
	if err != nil {
		return
	}

	go func() {
		fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()

	channel, err := conn.Channel()
	if err != nil {
		return
	}

	err = channel.ExchangeDeclare(
		"gusha",  // name of the exchange
		"direct", // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)

	if err != nil {
		return
	}

	queue, err := channel.QueueDeclare(
		"incoming", // name of the queue
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		return
	}

	h, _ := os.Hostname()
	err = channel.QueueBind(
		queue.Name, // name of the queue
		h,          // bindingKey
		"gusha",    // sourceExchange
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		return
	}

	deliveries, err := channel.Consume(
		queue.Name, // name
		h,          // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		return
	}

	go handle(deliveries, done)

	return nil
}
