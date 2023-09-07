package scheduler

import (
	"fmt"
	"log"
	"test/internal/config"
	"test/internal/entity"
	"test/internal/messager"
	"test/internal/metrics"
	"test/internal/usecase"
	"time"

	"github.com/roylee0704/gron"
)

type Scheduler struct {
	uc        usecase.Usecase
	cfg       config.Configurator
	gron      *gron.Cron
	processor *messager.Processor
}

func NewScheduler(uc usecase.Usecase, cfg config.Configurator, processor *messager.Processor) *Scheduler {
	//create scheduler object
	s := &Scheduler{
		uc:        uc,
		cfg:       cfg,
		gron:      gron.New(),
		processor: processor,
	}

	// add workers
	workers := cfg.GetWorkerHandlers()
	for _, worker := range workers {
		if err := s.AddWorker(worker); err != nil {
			log.Fatal("scheduler worker add error: ", err)
		}
	}

	return s
}

func (s *Scheduler) Run() {
	s.gron.Start()
}

func (s *Scheduler) Stop() {
	s.gron.Stop()
}

func (s *Scheduler) AddWorker(worker entity.Worker) error {
	period, err := time.ParseDuration(worker.Period)
	if err != nil {
		return fmt.Errorf("failed add worker: %w", err)
	}

	schedule := gron.Every(period)
	s.gron.AddFunc(schedule, func() {

		//debug, delete it after test
		testMap := map[int]string{
			0: "",
			1: "test",
		}
		values := []string{}
		for _, v := range testMap {
			values = append(values, v)
		}
		//

		log.Printf("Run worker: %s\n", worker.Name)
		outputs, err := worker.HandlerFn(worker)
		if err != nil {
			log.Printf("worker function recieved error: %s\n", err.Error())

			return
		}

		log.Printf("Send worker outputs: %s\n", worker.Name)
		err = s.processor.HandleWorker(outputs, worker)
		if err != nil {
			log.Printf("handle worker error: %s\n", err.Error())
			metrics.WordsOperationResults.WithLabelValues(values[0], values[1]).Set(0)
			return
		}
		metrics.WordsOperationResults.WithLabelValues(values[0], values[1]).Set(1)

		log.Printf("Worker `%s` successfully finished\n", worker.Name)
	})

	log.Printf("Scheduler worker `%s` added (%s)\n", worker.Name, worker.Period)

	return nil
}

var cnt int

func Do(fn func()) {
	if cnt == 2 {
		return
	}
	fn()
	cnt++
}
