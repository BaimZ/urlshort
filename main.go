package main

import (
	"context"
	"flag"
	"log"

	tgClient "urlshortener/clients/telegram"
	event_consumer "urlshortener/consumer/event-consumer"
	"urlshortener/events/telegram"
	"urlshortener/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("cant connect to storage: %s", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to tgbot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}
	return *token
}
