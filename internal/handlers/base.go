package handlers

import (
	"fmt"
	"test/internal/entity"
	"test/internal/usecase"

	"context"
)

func (h *Handlers) Base(input entity.Input) entity.Output {
	answer := "Welcome to phrases learn bot. Write a phrase for the reminder."
	data := input.GetData()

	out := input.CreateOutput()

	switch {
	case data.Type == entity.DataTypeMsg && data.Content != "/start":
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		if err := h.uc.CreatePhrase(ctx, input.GetUserID(), data.Content); err != nil {
			out.SetError(fmt.Errorf("save phrase error :%w", err))
		}

		answer = "phrase successfully added"

	case data.Type == entity.DataTypeMsg && data.Content == "/start":
		// create settings for user
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		h.uc.SetPhrasesNum(ctx, input.GetUserID(), usecase.DefaultPhrasesNum)
	}

	return out.
		SetMessage(answer).
		SetKeyboard(input.GetKeyboard()).
		SetUserID(input.GetUserID()).
		SetCache(input.GetCache())
}
