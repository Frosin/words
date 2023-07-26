package handlers

import (
	"fmt"
	"reflect"
	"test/internal/entity"
	"test/internal/service"
	"test/internal/usecase"
)

type Handlers struct {
	uc             usecase.Usecase
	serviceCfg     service.ServiceConfig
	handlers       map[string]entity.Handler
	workerHandlers map[string]entity.WorkerHandler
}

func NewHandlers(uc usecase.Usecase, serviceCfg service.ServiceConfig) *Handlers {
	handlers := &Handlers{
		uc:         uc,
		serviceCfg: serviceCfg,
	}

	handlers.handlers, handlers.workerHandlers = handlers.GetHandlers()

	return handlers
}

func (h *Handlers) GetHandlers() (map[string]entity.Handler, map[string]entity.WorkerHandler) {
	if h.handlers != nil && h.workerHandlers != nil {
		return h.handlers, h.workerHandlers
	}

	val := reflect.ValueOf(h)
	typ := reflect.TypeOf(h)

	result := make(map[string]entity.Handler, typ.NumMethod())
	resultWorkers := make(map[string]entity.WorkerHandler, typ.NumMethod())

	for i := 0; i < typ.NumMethod(); i++ {
		vMethod := val.Method(i)
		tMethod := typ.Method(i)

		switch hndlr := vMethod.Interface().(type) {
		case func(entity.Input) entity.Output:
			result[tMethod.Name] = hndlr
		case func(entity.Worker) ([]entity.Output, error):
			resultWorkers[tMethod.Name] = hndlr
		}
	}

	h.handlers, h.workerHandlers = result, resultWorkers

	return result, resultWorkers
}

func (h *Handlers) findHandler(name string) (entity.Handler, error) {
	handler, ok := h.handlers[name]
	if !ok {
		return nil, fmt.Errorf("handler not found: %s", name)
	}

	return handler, nil
}
