package messager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"test/internal/entity"
	"test/internal/metrics"
	"test/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func convertKeyboard(kbd entity.Keyboard) *tgbotapi.InlineKeyboardMarkup {
	if len(kbd.Buttons) == 0 {
		return nil
	}

	controlButtons := []tgbotapi.InlineKeyboardButton{}
	for _, btn := range kbd.Buttons {
		controlButtons = append(
			controlButtons,
			tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Handler))
	}

	result := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			controlButtons,
		},
	}

	return &result
}

func newData(msg, cmd *string) entity.Data {
	switch {
	case msg != nil:
		return entity.Data{
			Content: *msg,
			Type:    entity.DataTypeMsg,
		}
	case cmd != nil:
		return entity.Data{
			Content: *cmd,
			Type:    entity.DataTypeCmd,
		}
	default:
		panic("newData unexpected behavior")
	}
}

func (p *Processor) handleUpdate(update tgbotapi.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.sc.DBTimeout())
	defer cancel()

	var (
		fromUser *tgbotapi.User
		msg      *string
		cmd      *string
		chatID   int64
	)
	switch {
	case update.CallbackQuery != nil:
		fromUser = update.CallbackQuery.From
		cmd = &update.CallbackQuery.Data
		chatID = update.CallbackQuery.Message.Chat.ID
	case update.Message != nil:
		fromUser = update.Message.From
		msg = &update.Message.Text
		chatID = update.Message.Chat.ID
	default:
		return fmt.Errorf("unexpected error: no data")
	}

	// get session by updateData
	session, err := p.repo.GetSession(ctx, int64(fromUser.ID))
	switch {
	case session == nil && err == nil:
		session = &entity.Session{
			CurrentPage: p.cfg.GetFirstPageName(),
			UserID:      fromUser.ID,
			LastMsgID:   0,
			Data:        "",
		}
	case err != nil:
		return fmt.Errorf("get session error: %w", err)
	}

	// check for first settings
	settings, err := p.repo.GetSettings(ctx, int64(fromUser.ID))
	switch {
	case settings == nil && err == nil:
		settings := entity.UserSettings{
			UserID:         int64(fromUser.ID),
			PhraseDayLimit: usecase.DefaultPhrasesNum,
		}
		if err := p.repo.SaveSettings(ctx, settings); err != nil {
			return fmt.Errorf("save first settings error: %w", err)
		}
	case err != nil:
		return fmt.Errorf("get settings error: %w", err)
	}

	sessData, err := session.Data.GetSessionData()
	if err != nil {
		return fmt.Errorf("get session data error: %w", err)
	}

	curPage := p.cfg.FindPage(session.CurrentPage)

	data := newData(msg, cmd)

	// find handler and page name for calculating
	handler, handlerPage, err := p.cfg.FindHandler(curPage, data)
	if err != nil {
		return err
	}

	input := entity.NewInput(
		data,
		handlerPage.StartKeyboard,
		int64(fromUser.ID),
		sessData,
		session.CurrentPhraseNum,
		session.CurrentDay,
	)

	output := handler(input)

	if err := output.GetError(); err != nil {
		err := fmt.Errorf("handler error on page %s with data (msg=%v, data=%v) :%w", curPage.Name, msg, data, output.GetError())
		log.Println("handler error", err)
		return err
	}

	rawSessData, err := output.GetCache().ToRaw()
	if err != nil {
		return fmt.Errorf("serialize session data error: %w", err)
	}

	keyboard := convertKeyboard(output.GetKeyboard())

	var response tgbotapi.Chattable
	//if session.LastMsgID == 0 {
	resp := tgbotapi.NewMessage(chatID, output.GetMessage())
	if keyboard != nil {
		resp.ReplyMarkup = *keyboard
	}
	response = resp
	// } else {
	// 	resp := tgbotapi.NewEditMessageText(
	// 		chatID,
	// 		session.LastMsgID,
	// 		output.GetMessage(),
	// 	)
	// 	if keyboard != nil {
	// 		resp.ReplyMarkup = keyboard
	// 	}
	// 	response = resp
	// }

	// delete keyboard from last message
	if session.LastMsgHasKbd {
		previous := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			session.LastMsgID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		_, err = p.bot.Send(previous)
		if err != nil {
			log.Println("error sending clear previous msg")

			return err
		}
	}

	// send new message
	mess, err := p.bot.Send(response)
	if err != nil {
		log.Println("error sending msg")

		return err
	}

	if output.GetGoToStart() {
		session.CurrentPage = p.cfg.GetFirstPageName()
	} else {
		// save current page name
		session.CurrentPage = handlerPage.Name
	}

	session.LastMsgID = mess.MessageID

	session.ChatID = chatID
	session.Data = rawSessData

	if output.GetCurrentDay() != 0 {
		session.CurrentPhraseNum = output.GetCurrentPhraseNum()
		session.CurrentDay = output.GetCurrentDay()
	}

	if keyboard != nil {
		session.LastMsgHasKbd = true
	} else {
		session.LastMsgHasKbd = false
	}

	// save session
	if err := p.repo.SaveSession(ctx, *session); err != nil {
		return err
	}

	return nil
}

func (p *Processor) RegisterAndRunTelegramBot() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := p.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	//debug, delete it after test
	testMap := map[int]string{
		0: "",
		1: "test",
	}
	values := []string{}
	for _, v := range testMap {
		values = append(values, v)
	}
	//

	for update := range updates {
		update := update
		go func() {
			if err := p.handleUpdate(update); err != nil {
				metrics.WordsOperationResults.WithLabelValues(values[0], values[1]).Set(0)
				log.Println("ERROR", err)
			} else {
				metrics.WordsOperationResults.WithLabelValues(values[0], values[1]).Set(1)
			}
		}()
	}

	return errors.New("unexpected exit")
}
