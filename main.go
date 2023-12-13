package main

import (
	"flag"
	"log"
	"urlshortener/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

	//fetcher = fetcher.New()
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to tgbot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}

}
