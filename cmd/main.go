package main

import (
	"log"
	"test/internal/config"
	"test/internal/messager"
	"test/internal/metrics"
	"test/internal/repository"
	"test/internal/scheduler"
	"test/internal/service"
	"test/internal/usecase"
)

func main() {
	serviceConfig := service.NewServiceConfig("")

	repo := repository.NewRepository(serviceConfig)

	usecase := usecase.NewUsecase(repo)

	botConfig := config.NewBotConfig(serviceConfig, usecase)

	botProcessor := messager.NewProcessor(botConfig, repo, serviceConfig)

	scheduler := scheduler.NewScheduler(usecase, botConfig, botProcessor)
	scheduler.Run()

	metrics.RunMetrics()

	err := botProcessor.RegisterAndRunTelegramBot()
	if err != nil {
		log.Fatal(err)
	}
}
