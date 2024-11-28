package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sebastianaldi17/ngormq/internal/entity"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Handler struct {
	db      *gorm.DB
	channel *amqp.Channel
	queue   amqp.Queue
}

func New(channel *amqp.Channel, db *gorm.DB) (*Handler, error) {
	handler := &Handler{}
	if channel != nil {
		q, err := channel.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			return nil, err
		}
		handler.channel = channel
		handler.queue = q
	}

	if db != nil {
		handler.db = db
	}

	return handler, nil
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
		http.Error(w, "could not get hostname", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Pong from %s\n", hostname)
}

func (h *Handler) ConsumeNewMessage(msg amqp.Delivery) error {
	var message entity.Message
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	messageModel := entity.Message{
		Hostname: hostname,
		Message:  "Hello world!",
	}
	if h.db != nil {
		res := h.db.Create(&messageModel)
		return res.Error
	}

	return errors.New("DB not initialized")
}

func (h *Handler) PublishNewMessage() error {
	msg, err := json.Marshal(entity.Message{
		Message: "Hello world!",
	})
	if err != nil {
		return err
	}

	// get from https://api.quotable.io/random

	err = h.channel.Publish(
		"",           // exchange
		h.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	return err
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB not initialized"))
		return
	}

	messages := make([]entity.Message, 0)
	res := h.db.Find(&messages)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res.Error.Error()))
		return
	}

	out, _ := json.MarshalIndent(messages, "", "  ")
	w.Write(out)
}
