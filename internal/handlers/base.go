package handlers

import (
	"fmt"
	"test/internal/entity"
	"test/internal/metrics"

	"context"
)

func (h *Handlers) Base(input entity.Input) entity.Output {
	answer := "Welcome to phrases learn bot. Write a phrase for the reminder."
	data := input.GetData()

	out := input.CreateOutput()

	if data.Type == entity.DataTypeMsg && data.Content != "/start" {
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		if err := h.uc.CreatePhrase(ctx, input.GetUserID(), data.Content); err != nil {
			out.SetError(fmt.Errorf("save phrase error :%w", err))
		}

		// add phrase to cache
		cache := make(entity.SessionData, 1)
		cache["phrase"] = data.Content
		out.SetCache(cache)

		answer = "phrase successfully added"
		// send metric
		metrics.WordsPhraseAdded.Inc()
	} else {
		out.SetCache(entity.NewSessionData())
	}

	return out.
		SetMessage(answer).
		SetUserID(input.GetUserID()).
		SetKeyboard(input.GetKeyboard())
}
