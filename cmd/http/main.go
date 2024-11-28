package main

import (
	"log"
	"net/http"

	"github.com/sebastianaldi17/ngormq/internal/handler"

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

	handler, err := handler.New(nil, db)
	if err != nil {
		log.Panicf("%s", err)
	}

	http.HandleFunc("/ping", handler.Ping)
	http.HandleFunc("/messages", handler.GetMessages)
	log.Println("Starting HTTP server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
