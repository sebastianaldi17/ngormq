package main

import (
	"log"
	"time"

	"github.com/sebastianaldi17/ngormq/internal/handler"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Panicf("%s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s", err)
	}
	defer ch.Close()

	handler, err := handler.New(ch, nil)
	if err != nil {
		log.Panicf("%s", err)
	}

	log.Println("Begin publishing messages")

	for {
		err = handler.PublishNewMessage()
		if err != nil {
			log.Printf("Error publishing new message: %s", err.Error())
		}
		time.Sleep(time.Second * 5)
	}
}
