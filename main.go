package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/streadway/amqp"
)

var (
	amqpURI = flag.String("rabbit", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func main() {
	flag.Parse()

	Consumer()

	for m := range done {
		log.Print(m)
	}
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		i := Incoming{}

		err := json.Unmarshal(d.Body, &i)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Starting %q", i.Name)

		agent, err := NewAgent(i.BaseURL, i.URLs, i.Seconds, i.Clients)
		if err != nil {
			log.Panic(err)
		}

		agent.Run()

		d.Ack(true)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
