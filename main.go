package main

import (
	"flag"
	"log"

	tgClient "tg-bot/clients/telegram"
	event_consumer "tg-bot/consumer/event-consumer"
	"tg-bot/events/telegram"
	"tg-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
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
		log.Fatal("service stopped with error", err)
	}

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
