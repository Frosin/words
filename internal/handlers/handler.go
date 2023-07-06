package handlers

import (
	"reflect"
	"test/internal/entity"
	"test/internal/service"
	"test/internal/usecase"
)

type Handlers struct {
	uc         usecase.Usecase
	serviceCfg service.ServiceConfig
}

func NewHandlers(uc usecase.Usecase, serviceCfg service.ServiceConfig) *Handlers {
	return &Handlers{
		uc:         uc,
		serviceCfg: serviceCfg,
	}
}

func (h *Handlers) GetHandlers() (map[string]entity.Handler, map[string]entity.WorkerHandler) {
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

	return result, resultWorkers
}
