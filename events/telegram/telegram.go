package telegram

import (
	"tg-bot/clients/telegram"
	"tg-bot/events"
	"tg-bot/lib/errorr"
	"tg-bot/storage"
)

type Processor struct {
	tg *telegram.Client
	offset int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg: client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errorr.Wrap("can't get events", err)
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}
}

func event(update telegram.Update) events.Event {
	res := events.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}
}