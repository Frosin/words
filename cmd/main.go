package main

import (
	"log"
	"test/internal/backup"
	"test/internal/config"
	"test/internal/messager"
	"test/internal/metrics"
	"test/internal/repository"
	"test/internal/scheduler"
	"test/internal/service"
	"test/internal/usecase"
)

func createDumpFunction(sc *service.SConfig) backup.DumpFn {
	yaDiskToken := sc.GetYadiskToken()
	dbPath := sc.GetDBFileName()

	return func() error {
		return backup.UploadBackupDB(yaDiskToken, "DB", dbPath, true)
	}
}

func main() {
	serviceConfig := service.NewServiceConfig("")

	repo := repository.NewRepository(serviceConfig)

	dumper := backup.NewDumper(createDumpFunction(serviceConfig), nil)
	dumper.Start()

	usecase := usecase.NewUsecase(repo, dumper)

	botConfig := config.NewBotConfig(serviceConfig, usecase)

	botProcessor := messager.NewProcessor(botConfig, repo, serviceConfig)

	scheduler := scheduler.NewScheduler(usecase, botConfig, botProcessor, serviceConfig.IsWorkerDebugOnce())
	scheduler.Run()

	metrics.RunMetrics()

	err := botProcessor.RegisterAndRunTelegramBot()
	if err != nil {
		log.Fatal(err)
	}
}
