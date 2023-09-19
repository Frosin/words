package messager

import (
	"context"
	"fmt"
	"log"
	"test/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MessagerWorker interface {
	HandleWorker(out entity.Output) error
}

func (p *Processor) HandleWorker(outs []entity.Output, worker entity.Worker) []error {
	ctx, cancel := context.WithTimeout(context.Background(), p.sc.DBTimeout())
	defer cancel()

	if len(outs) == 0 {
		return nil
	}

	// set one output = one user
	// first output will be output with the earliest created_at and the least epoch
	userIDs := make([]int64, 0, len(outs))
	outMap := make(map[int64]entity.Output, len(outs))
	for _, out := range outs {
		outMap[out.GetUserID()] = out
	}

	for userID := range outMap {
		userIDs = append(userIDs, userID)
	}

	// get session by updateData
	sessions, err := p.repo.GetSessions(ctx, userIDs)
	switch {
	case len(sessions) == 0 && err == nil:
		log.Printf("worker failed: not found session for users = %v\n", userIDs)

		return nil
	case err != nil:
		return []error{fmt.Errorf("worker get session error: %w", err)}
	}

	// find action handler and page name for calculating
	workerPage := p.cfg.FindPage(worker.Page)
	if workerPage == nil {
		return []error{fmt.Errorf("worker page not found %s", worker.Page)}
	}

	errs := []error{}

	for _, session := range sessions {
		output, ok := outMap[int64(session.UserID)]
		if !ok {
			err := fmt.Errorf("not found output for user: %d", session.UserID)
			log.Println("get output error", err)

			errs = append(errs, err)
			continue
		}

		if err := output.GetError(); err != nil {
			err := fmt.Errorf("worker error with data: %w", output.GetError())
			log.Println("worker handler error", err)

			errs = append(errs, err)
			continue
		}

		keyboard := convertKeyboard(output.GetKeyboard())

		response := tgbotapi.NewMessage(session.ChatID, output.GetMessage())
		if keyboard != nil {
			response.ReplyMarkup = *keyboard
		}

		// delete keyboard from last message
		if session.LastMsgHasKbd {
			previous := tgbotapi.NewEditMessageReplyMarkup(
				session.ChatID,
				session.LastMsgID,
				tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
				},
			)
			_, err = p.bot.Send(previous)
			if err != nil {
				log.Printf("error sending clear previous msg: %v: %s\n", session, err.Error())

				errs = append(errs, err)
				continue
			}
		}

		// send new message
		mess, err := p.bot.Send(response)
		if err != nil {
			log.Printf("error sending msg: %s\n", err.Error())

			errs = append(errs, err)
			continue
		}

		rawSessData, err := output.GetCache().ToRaw()
		if err != nil {
			err = fmt.Errorf("serialize session data error: %w", err)

			errs = append(errs, err)
			continue
		}

		session.Data = rawSessData
		session.LastMsgID = mess.MessageID
		// save current page name
		session.CurrentPage = workerPage.Name

		if keyboard != nil {
			session.LastMsgHasKbd = true
		} else {
			session.LastMsgHasKbd = false
		}

		// save session
		if err := p.repo.SaveSession(ctx, *session); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errs
}
