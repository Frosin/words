package config

import (
	"fmt"
	"test/internal/entity"
)

func findHandler(handlers map[string]entity.Handler, handlerName string) (entity.Handler, error) {
	handlerFn, ok := handlers[handlerName]
	if !ok {
		return nil, fmt.Errorf("implementation of handler %s not found", handlerName)
	}

	return handlerFn, nil
}

func findWorkerHandler(handlers map[string]entity.WorkerHandler, handlerName string) (entity.WorkerHandler, error) {
	handlerFn, ok := handlers[handlerName]
	if !ok {
		return nil, fmt.Errorf("implementation of handler worker %s not found", handlerName)
	}

	return handlerFn, nil
}

func parse(cfg *entity.Config, handlers map[string]entity.Handler, workerHandlers map[string]entity.WorkerHandler) (map[string]*entity.Page, error) {
	pages := make(map[string]*entity.Page, len(cfg.Pages))
	for i := range cfg.Pages {
		pages[cfg.Pages[i].Name] = &cfg.Pages[i]
	}

	for i := range cfg.Pages {
		// find page handler
		handlerFn, err := findHandler(handlers, cfg.Pages[i].Handler)
		if err != nil {
			return nil, err
		}
		cfg.Pages[i].HandlerFn = handlerFn

		for j := range cfg.Pages[i].StartKeyboard.Buttons {
			btn := cfg.Pages[i].StartKeyboard.Buttons[j]
			handlerFn, err := findHandler(handlers, btn.Handler)
			if err != nil {
				return nil, err
			}
			cfg.Pages[i].StartKeyboard.Buttons[j].HandlerFn = handlerFn
		}

		// check for first page
		if cfg.Pages[i].First {
			cfg.FirstPage = &cfg.Pages[i]
		}
	}

	// parse workers
	for i := range cfg.Workers {
		workerHandlerFn, err := findWorkerHandler(workerHandlers, cfg.Workers[i].WorkerHandler)
		if err != nil {
			return nil, err
		}

		cfg.Workers[i].HandlerFn = workerHandlerFn

		// find a handler page
		if cfg.Workers[i].Page != "" {
			workerPageEntity := findPage(pages, cfg.Workers[i].Page)
			if workerPageEntity == nil {
				return nil, fmt.Errorf("worker page entity not found %s", cfg.Workers[i].Page)
			}

			cfg.Workers[i].PageEntity = workerPageEntity
		}
	}

	if cfg.FirstPage == nil {
		return nil, fmt.Errorf("first page flag not found")
	}

	return pages, nil
}
