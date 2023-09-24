package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"test/internal/entity"
	"test/internal/handlers"
	"test/internal/service"
	"test/internal/usecase"
)

func NewBotConfig(sc service.ServiceConfig, uc usecase.Usecase) Configurator {
	botConfigFile := sc.GetBotConfigFileName()

	if _, err := os.Stat(botConfigFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal("bot config file not found")
	}

	cfgJson, err := os.ReadFile(botConfigFile)
	if err != nil {
		log.Fatal("bot config file read error")
	}

	handlers := handlers.NewHandlers(uc, sc)

	handlerFns, workerHandlerFns := handlers.GetHandlers()

	cfg, pages, err := parseConfig(cfgJson, handlerFns, workerHandlerFns)
	if err != nil {
		log.Fatal("bot parse config error")
	}

	botCfg := newConfig(cfg, pages, handlerFns, workerHandlerFns)

	return botCfg
}

func parseConfig(cfgJson []byte, handlers map[string]entity.Handler, workerHandlers map[string]entity.WorkerHandler) (*entity.Config, map[string]*entity.Page, error) {
	cfg := entity.Config{}

	if err := json.Unmarshal(cfgJson, &cfg); err != nil {
		return nil, nil, err
	}

	pages, err := parse(&cfg, handlers, workerHandlers)
	if err != nil {
		return nil, nil, err
	}

	return &cfg, pages, nil
}
