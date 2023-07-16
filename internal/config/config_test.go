package config

import (
	"os"
	"test/internal/handlers"
	"test/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfig(t *testing.T) {
	cfgJson, err := os.ReadFile("example.json")
	assert.NoError(t, err)

	handlers := handlers.NewHandlers(nil, &service.SConfig{})

	handlerFns, workerHandlerFns := handlers.GetHandlers()

	_, _, err = parseConfig(cfgJson, handlerFns, workerHandlerFns)
	assert.NoError(t, err)
}
