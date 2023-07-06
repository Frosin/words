package handlers

import (
	"context"
	"fmt"
	"test/internal/entity"
)

func (h *Handlers) Settings(input entity.Input) entity.Output {
	out := input.CreateOutput()

	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	phrasesNum, err := h.uc.GetPhrasesNum(ctx, input.GetUserID())
	if err != nil {
		out.SetError(fmt.Errorf("get phrases num error :%w", err))
	}

	msg := fmt.Sprintf("current phrases number: %d", phrasesNum)

	return out.
		SetMessage(msg).
		SetKeyboard(input.GetKeyboard()).
		SetUserID(input.GetUserID()).
		SetCache(input.GetCache())
}

func (h *Handlers) Settings_up(input entity.Input) entity.Output {
	out := input.CreateOutput()

	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	phrasesNum, err := h.uc.GetPhrasesNum(ctx, input.GetUserID())
	if err != nil {
		out.SetError(fmt.Errorf("get phrases num error :%w", err))
	}

	phrasesNum = phrasesNum + 1

	if err := h.uc.SetPhrasesNum(ctx, input.GetUserID(), phrasesNum); err != nil {
		out.SetError(fmt.Errorf("set phrases num error :%w", err))
	}

	msg := fmt.Sprintf("current phrases number: %d", phrasesNum)

	return out.
		SetMessage(msg).
		SetKeyboard(input.GetKeyboard()).
		SetUserID(input.GetUserID()).
		SetCache(input.GetCache())
}

func (h *Handlers) Settings_down(input entity.Input) entity.Output {
	out := input.CreateOutput()

	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	phrasesNum, err := h.uc.GetPhrasesNum(ctx, input.GetUserID())
	if err != nil {
		out.SetError(fmt.Errorf("get phrases num error :%w", err))
	}

	phrasesNum = phrasesNum - 1

	if err := h.uc.SetPhrasesNum(ctx, input.GetUserID(), phrasesNum); err != nil {
		out.SetError(fmt.Errorf("set phrases num error :%w", err))
	}

	msg := fmt.Sprintf("current phrases number: %d", phrasesNum)

	return out.
		SetMessage(msg).
		SetKeyboard(input.GetKeyboard()).
		SetUserID(input.GetUserID()).
		SetCache(input.GetCache())
}
