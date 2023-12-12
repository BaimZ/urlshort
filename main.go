package main

import (
	"flag"
	"log"
)

func main() {
	tgClient = telegram.New(token)

	fetcher = fetcher.New()
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
