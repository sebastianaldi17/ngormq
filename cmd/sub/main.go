package main

import (
	"log"

	"github.com/sebastianaldi17/ngormq/internal/handler"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(
		postgres.Open(`host=postgres user=username password=password dbname=sample-app port=5432 sslmode=disable TimeZone=Asia/Jakarta`),
		&gorm.Config{},
	)
	if err != nil {
		log.Panicf("%s", err)
	}

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

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("%s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("%s", err)
	}

	handler, err := handler.New(nil, db)
	if err != nil {
		log.Panicf("%s", err)
	}

	var block chan struct{}

	log.Println("Begin consuming messages")

	go func() {
		for d := range msgs {
			err := handler.ConsumeNewMessage(d)
			if err != nil {
				log.Printf("Error trying to consume message: %s, error: %s", string(d.Body), err.Error())
			}
		}
	}()

	<-block
}
