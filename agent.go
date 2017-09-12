package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

// Agent contains such contextual information for an agent
// such as configuration and rabbitmq pushers
type Agent struct {
	rabbit *amqp.Channel

	baseURL string
	urls    URLs
	seconds int
	clients int
}

// NewAgent spans an agent and, along with it, necessary rabbits
// and stuff
func NewAgent(baseURL string, urls URLs, seconds, clients int) (a Agent, err error) {
	a.baseURL = baseURL
	a.urls = urls
	a.seconds = seconds
	a.clients = clients

	connection, err := amqp.Dial(*amqpURI)
	if err != nil {
		return
	}

	a.rabbit, err = connection.Channel()
	if err != nil {
		return
	}

	err = a.rabbit.ExchangeDeclare(
		"gusha",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)

	return
}

// Run starts, runs, and handles a run through of loadtests
func (a Agent) Run() {
	BaseURL = a.baseURL

	channels := make([]chan bool, a.clients)

	for idx := 0; idx < a.clients; idx++ {
		log.Printf("Starting worker %d", idx)

		c := make(chan bool)
		go a.worker(c)

		channels = append(channels, c)
	}

	timer1 := time.NewTimer(time.Second * time.Duration(a.seconds))
	<-timer1.C

	log.Printf("%d seconds complete!", a.seconds)

	for idx, c := range channels {
		log.Printf("Completing worker %d", idx)
		c <- true
	}
}

func (a Agent) worker(c chan bool) {
	rand.Seed(time.Now().Unix())

	for {
		u := a.urls[rand.Intn(len(a.urls))]
		response, err := u.Do()
		if err != nil {
			log.Print(err)
		}

		go a.syncUp(response)
	}
}

func (a Agent) syncUp(r *http.Response) {
	body, err := json.Marshal(map[string]interface{}{
		"Status": r.Status,
		"URL":    r.Request.URL.String(),
		"Size":   r.ContentLength,
	})
	if err != nil {
		log.Print(err)
	}

	err = a.rabbit.Publish(
		"gusha", // publish to an exchange
		"gusha", // routing to 0 or more queues
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)

	if err != nil {
		log.Panic(err)
	}
}
