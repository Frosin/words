package handlers

import (
	"encoding/base64"
	"fmt"
	"test/internal/entity"
	"test/internal/metrics"

	"context"
)

var (
	phraseKey   = "phrase"
	sentenceKey = "sentence"
)

func (h *Handlers) Base(input entity.Input) entity.Output {

	data := input.GetData()

	out := input.CreateOutput()

	kbd := input.GetKeyboard()

	decodedAnswer, err := base64.StdEncoding.DecodeString(`0JTQvtCx0YDQviDQv9C+0LbQsNC70L7QstCw0YLRjCDQsiDQsdC+0YIg0LTQu9GPINC30LDQv9C+0LzQuNC90LDQvdC40Y8g0YHQu9C+0LIuINCg0LDQsdC+0YLQsNC10YIg0L/QviDQv9GA0LjQvdGG0LjQv9GDIFvQutGA0LjQstC+0Lkg0LfQsNCx0YvQstCw0L3QuNGPXShodHRwczovL3J1Lndpa2lwZWRpYS5vcmcvd2lraS8lRDAlOUElRDElODAlRDAlQjglRDAlQjIlRDAlQjAlRDElOEZfJUQwJUI3JUQwJUIwJUQwJUIxJUQxJThCJUQwJUIyJUQwJUIwJUQwJUJEJUQwJUI4JUQxJThGKQ0K0J/RgNC+0YHRgtC+INC/0LjRiNC40YLQtSDRgdGO0LTQsCDRgdC70L7QstCwINC40LvQuCDQstGL0YDQsNC20LXQvdC40Y8g0LrQvtGC0L7RgNGL0LUg0YXQvtGC0LjRgtC1INC30LDQv9C+0LzQvdC40YLRjCwg0LHQvtGCINCx0YPQtNC10YIg0L/QtdGA0LjQvtC00LjRh9C10YHQutC4INC90LDQv9C+0LzQuNC90LDRgtGMINC+INC90LjRhSDQuCDQv9GA0L7RgdC40YLRjCDQstCy0LXRgdGC0Lgg0L/RgNC10LTQu9C+0LbQtdC90LjQtSDRgdC+0LTQtdGA0LbQsNGJ0LXQtSDQstCy0LXQtNC10L3QvdC+INCy0LDQvNC40YHQu9C+0LLQviDQuNC70Lgg0LLRi9GA0LDQttC10L3QuNC1Lg0K0JrQvtC70LjRh9C10YHRgtCy0L4g0YHQu9C+0LIg0LTQu9GPINC90LDQv9C+0LzQuNC90LDQvdC40Y8g0LIg0LTQtdC90Ywg0L3QsNGB0YLRgNCw0LjQstCw0LXRgtGB0Y8sINCx0L7RgiDQvdCw0YXQvtC00LjRgtGB0Y8g0LIg0LDQutGC0LjQstC90L7QuSDRgNCw0LfRgNCw0LHQvtGC0LrQtSwg0L/RgNC10LTQu9C+0LbQtdC90LjRjyDQuCDQvtGI0LjQsdC60Lgg0LIg0YDQsNCx0L7RgtC1INC80L7QttC90L4g0L3QsNC/0YDQsNCy0LvRj9GC0YwgQEhlbGVuX2ZyYW5jaCA=`)
	if err != nil {
		out.SetError(fmt.Errorf("decode answer error :%w", err))
		return out
	}
	answer := string(decodedAnswer)

	if data.Type == entity.DataTypeMsg && data.Content != "/start" {
		ctx, cancel := context.WithTimeout(context.Background(), h.serviceCfg.DBTimeout())
		defer cancel()

		if err := h.uc.CreatePhrase(ctx, input.GetUserID(), data.Content); err != nil {
			out.SetError(fmt.Errorf("save phrase error :%w", err))
		}

		// add phrase to cache
		cache := make(entity.SessionData, 1)
		cache[phraseKey] = data.Content
		out.SetCache(cache)

		answer = "phrase successfully added"
		// send metric
		metrics.WordsPhraseAdded.Inc()

		// schedule the backup
		h.uc.ScheduleBackUp()
	} else {
		out.SetCache(entity.NewSessionData())

		// if we do not have phrase we should remove delete button
		cache := input.GetCache()
		if _, ok := cache[phraseKey]; !ok {
			if len(kbd.Buttons) != 2 {
				out.SetError(fmt.Errorf("unexpected length of buttons: <> 2"))

				return out
			}

			kbd.Buttons = kbd.Buttons[:1]
		}
	}

	return out.
		SetMessage(answer).
		SetUserID(input.GetUserID()).
		SetKeyboard(kbd)
}
