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

		sentences, err := getSentences(uph)
		if err != nil {
			return nil, fmt.Errorf("GetReminderPhrases getSentences error: %w", err)
		}

		msg := fmt.Sprintf(remindMsgtemplate, uph.Phrase)

		if len(sentences) > 0 {
			msg += "\nУже сохраненные предложения:\n"
			for i, sen := range sentences {
				msg += fmt.Sprintf("%d. %s\n", i+1, sen)
			}
		}

		kbd := worker.PageEntity.StartKeyboard
		// if len(kbd.Buttons) != 1 {
		// 	return nil, fmt.Errorf("reminder keyboard has invalid num of btns")
		// }

		// // add cmd(phrase) to handler
		// kbd.Buttons[0].Handler = entity.CreateMsgCmdContent(kbd.Buttons[0].Handler, uph.Phrase)

		cache := make(entity.SessionData, 1)
		cache[phraseKey] = uph.Phrase

		out.SetMessage(msg).
			SetKeyboard(kbd).
			SetUserID(uph.UserID).
			SetCache(cache)

		outs = append(outs, out)
	}

	return outs, nil
}

func getSentences(ph *entity.Phrase) ([]string, error) {
	if ph == nil {
		return nil, nil
	}

	meta, err := entity.DeserializePhraseMeta(ph.Meta)
	if err != nil {
		return nil, err
	}

	return meta.Sentences, nil
}

func (h *Handlers) Page_reminder(input entity.Input) entity.Output {
	data := input.GetData()

	out := input.CreateOutput()

	cache := input.GetCache()

	phraseIntf, ok := cache[phraseKey]
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
	// if len(kbd.Buttons) != 1 {
	// 	out.SetError(fmt.Errorf("page reminder keyboard has invalid num of btns"))

	// 	return out
	// }

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

		err := h.uc.UpdatePhrase(ctx, input.GetUserID(), phrase, data.Content)
		if err != nil {
			out.SetError(fmt.Errorf("update phrase info :%w", err))
			return out
		}
		out.SetMessage("Предложение сохранено")

		// add sentence to cache
		cache[sentenceKey] = data.Content

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

		// add settings and delete buttons to keyboard
		settingsBtnHandler, err := h.findHandler("Settings")
		if err != nil {
			out.SetError(err)
			return out
		}

		deleteBtnHandler, err := h.findHandler("Delete_sentence")
		if err != nil {
			out.SetError(err)
			return out
		}

		newBtns := []entity.Button{
			{
				Text:      "Settings",
				Handler:   "Settings",
				HandlerFn: settingsBtnHandler,
			},
			{
				Text:      "❌ delete",
				Handler:   "Delete_sentence",
				HandlerFn: deleteBtnHandler,
			},
		}

		kbd.Buttons = append(kbd.Buttons, newBtns...)

		// go to first page next time
		out.SetGoToStart()

	case entity.DataTypeCmd:
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
		SetCache(cache)
}

func (h *Handlers) Delete_sentence(input entity.Input) entity.Output {
	out := input.CreateOutput()

	ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
	defer cancel()

	cache := input.GetCache()

	phraseIntf, ok := cache[phraseKey]
	if !ok {
		out.SetError(fmt.Errorf("phrase not found in cache"))
		return out
	}

	phrase, ok := phraseIntf.(string)
	if !ok {
		out.SetError(fmt.Errorf("phrase invalid type"))
		return out
	}

	sentenceIntf, ok := cache[sentenceKey]
	if !ok {
		out.SetError(fmt.Errorf("sentence not found in cache"))
		return out
	}

	sentence, ok := sentenceIntf.(string)
	if !ok {
		out.SetError(fmt.Errorf("sentence invalid type"))
		return out
	}

	err := h.uc.DeletePhraseSentence(ctx, input.GetUserID(), phrase, sentence)
	if err != nil {
		out.SetError(fmt.Errorf("delete phrases error :%w", err))
	}

	msg := fmt.Sprintf("sentence '%s' deleted", sentence)

	// remove delete button
	// because we already deleted phrase or we didn't write it yet
	kbd := input.GetKeyboard()
	// if len(kbd.Buttons) != 2 {
	// 	out.SetError(fmt.Errorf("invalid num of buttons"))
	// 	return out
	// }

	// remove delete button
	kbd.Buttons = kbd.Buttons[:len(kbd.Buttons)-1]

	return out.
		SetMessage(msg).
		SetKeyboard(kbd).
		SetUserID(input.GetUserID())
}
