package handlers

import (
	"context"
	"fmt"
	"test/internal/entity"
	"time"
)

const remindMsgtemplate = `Напишите предложение содержащее "%s"`

func (h *Handlers) Reminder(worker entity.Worker) ([]entity.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	// clear current session limits
	err := h.uc.ClearUsersDayLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("ClearUsersDayLimits error :%w", err)
	}

	// phrases ordered by user
	userPhrases, err := h.uc.GetReminderPhrases(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetReminderPhrases error :%w", err)
	}

	outs := make([]entity.Output, 0, len(userPhrases))

	for _, uph := range userPhrases {
		out := entity.NewOutput()

		msg := fmt.Sprintf(remindMsgtemplate, uph.Phrase)

		// kbd := worker.PageEntity.StartKeyboard
		// if len(kbd.Buttons) != 1 {
		// 	return nil, fmt.Errorf("reminder keyboard has invalid num of btns")
		// }

		// // add cmd(phrase) to handler
		// kbd.Buttons[0].Handler = entity.CreateMsgCmdContent(kbd.Buttons[0].Handler, uph.Phrase)

		cache := make(entity.SessionData, 1)
		cache["phrase"] = uph.Phrase

		out.SetMessage(msg).
			//SetKeyboard(kbd).
			SetUserID(uph.UserID).
			SetCache(cache)

		outs = append(outs, out)
	}

	return outs, nil
}

func (h *Handlers) Page_reminder(input entity.Input) entity.Output {
	data := input.GetData()

	out := input.CreateOutput()

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

	kbd := input.GetKeyboard()
	if len(kbd.Buttons) != 1 {
		out.SetError(fmt.Errorf("page reminder keyboard has invalid num of btns"))

		return out
	}

	// add cmd(phrase) to handler
	//if !entity.HasCmd(kbd.Buttons[0].Handler) {
	//kbd.Buttons[0].Handler = entity.CreateMsgCmdContent(kbd.Buttons[0].Handler, phrase)
	//}

	switch data.Type {
	case entity.DataTypeMsg:
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		// TODO
		// put here request to chatGPT
		// and add to output message gpt's sentence

		h.uc.UpdatePhrase(ctx, input.GetUserID(), phrase, data.Content)
		out.SetMessage("Предложение сохранено")

		// phrases limit logic:
		// check session day
		// and increase current words num
		// or set it to zero
		curDay := time.Now().UTC().Day()

		sessCurDay := input.GetCurrentDay()
		sessCurWordsNum := input.GetCurrentPhraseNum()
		if curDay != sessCurDay {
			out.SetCurrentDay(curDay)
			out.SetCurrentPhraseNum(1)
		} else {
			out.SetCurrentDay(curDay)
			out.SetCurrentPhraseNum(sessCurWordsNum + 1)
		}

	case entity.DataTypeCmd:
		// create settings for user
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		info, err := h.uc.GetPhraseInfo(ctx, input.GetUserID(), phrase)
		if err != nil {
			out.SetError(fmt.Errorf("get phrase info :%w", err))
			return out
		}

		message := fmt.Sprintf(
			"Фраза: %s, добавлена: %s, этап запоминания: %d-й",
			info.Phrase,
			info.CreatedAt.Format("02.01.2006"),
			info.Epoch,
		)
		out.
			SetMessage(message)
	}

	return out.
		SetKeyboard(kbd).
		SetUserID(input.GetUserID()).
		SetGoToStart().
		SetCache(cache)
}
