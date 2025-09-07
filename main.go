package main

import (
	"flag"
	"log"
	"os"

	tgClient "tg-bot/clients/telegram"
	event_consumer "tg-bot/consumer/event-consumer"
	"tg-bot/events/telegram"
	"tg-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	// Создаем базовую директорию для storage перед использованием
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Can't create storage directory: %v", err)
	}

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
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not set")
	}

	return *token
}
