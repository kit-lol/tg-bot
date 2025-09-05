package main

import (
	"flag"
	"log"
	"tg-bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
}

func mustToken() string {
	token := flag.String(
		"t", 
		"", 
		"token for access ot telegram bot",
	)
	
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not set")
	}

	return *token
}