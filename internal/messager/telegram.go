package messager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"test/internal/entity"
	"test/internal/metrics"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func convertKeyboard(kbd entity.Keyboard) *tgbotapi.InlineKeyboardMarkup {
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

	for update := range updates {
		update := update
		go func() {
			if err := p.handleUpdate(update); err != nil {
				metrics.WordsOperationResults.Set(0)
				log.Println("ERROR", err)
			} else {
				metrics.WordsOperationResults.Set(1)
			}
		}()
	}

	return errors.New("unexpected exit")
}
