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

	phrasesNum, err := h.uc.GetPhrasesDayLimit(ctx, input.GetUserID())
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

	phrasesNum, err := h.uc.GetPhrasesDayLimit(ctx, input.GetUserID())
	if err != nil {
		out.SetError(fmt.Errorf("get phrases num error :%w", err))
	}

	phrasesNum = phrasesNum + 1

	if err := h.uc.SetPhrasesDayLimit(ctx, input.GetUserID(), phrasesNum); err != nil {
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

	phrasesNum, err := h.uc.GetPhrasesDayLimit(ctx, input.GetUserID())
	if err != nil {
		out.SetError(fmt.Errorf("get phrases num error :%w", err))
	}

	phrasesNum = phrasesNum - 1

	if err := h.uc.SetPhrasesDayLimit(ctx, input.GetUserID(), phrasesNum); err != nil {
		out.SetError(fmt.Errorf("set phrases num error :%w", err))
	}

	msg := fmt.Sprintf("current phrases number: %d", phrasesNum)

	return out.
		SetMessage(msg).
		SetKeyboard(input.GetKeyboard()).
		SetUserID(input.GetUserID()).
		SetCache(input.GetCache())
}

func (h *Handlers) Delete_phrase(input entity.Input) entity.Output {
	out := input.CreateOutput()

	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	cache := input.GetCache()

	phraseIntf, ok := cache["phrase"]
	if !ok {
		out.SetError(fmt.Errorf("phrase not found in cache"))
		return out
	}

	phrase, ok := phraseIntf.(string)
	if !ok {
		out.SetError(fmt.Errorf("phrase invalid type"))
		return out
	}

	err := h.uc.DeletePhrase(ctx, input.GetUserID(), phrase)
	if err != nil {
		out.SetError(fmt.Errorf("delete phrases error :%w", err))
	}

	msg := fmt.Sprintf("phrase '%s' deleted", phrase)

	// remove delete button
	// because we already deleted phrase or we didn't write it yet
	kbd := input.GetKeyboard()
	if len(kbd.Buttons) != 2 {
		out.SetError(fmt.Errorf("invalid num of buttons"))
		return out
	}
	kbd.Buttons = kbd.Buttons[:1]

	return out.
		SetMessage(msg).
		SetKeyboard(kbd).
		SetUserID(input.GetUserID())
}
