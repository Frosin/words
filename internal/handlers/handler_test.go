package handlers

import (
	"test/internal/repository"
	"test/internal/service"
	"test/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetHandlers(t *testing.T) {
	serviceCfg := service.SConfig{}
	repo := repository.NewRepository(&serviceCfg)
	uc := usecase.NewUsecase(repo, nil)
	h := NewHandlers(uc, &serviceCfg)

	handlers, workerhandlers := h.GetHandlers()
	assert.Len(t, handlers, 6)
	assert.Len(t, workerhandlers, 1)
}
