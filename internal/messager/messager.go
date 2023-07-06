package messager

import (
	"log"
	"test/internal/config"
	"test/internal/repository"
	"test/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Processor struct {
	cfg  config.Configurator
	bot  *tgbotapi.BotAPI
	repo repository.Repository
	sc   service.ServiceConfig
}

func NewProcessor(cfg config.Configurator, repo repository.Repository, sc service.ServiceConfig) *Processor {
	bot, err := tgbotapi.NewBotAPI(sc.GetBotToken())
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Processor{
		cfg, bot, repo, sc,
	}
}
