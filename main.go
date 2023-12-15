package main

import (
	"flag"
	"log"

	tgClient "urlshortener/clients/telegram"
	event_consumer "urlshortener/consumer/event-consumer"
	"urlshortener/events/telegram"

	"urlshortener/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
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
