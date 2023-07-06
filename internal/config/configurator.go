package config

import (
	"fmt"
	"test/internal/entity"
)

type Configurator interface {
	FindPage(name string) *entity.Page
	GetFirstPageName() string
	FindHandler(curPage *entity.Page, data entity.Data) (entity.Handler, *entity.Page, error)
	GetWorkerHandlers() []entity.Worker
}

type BotConfig struct {
	cfg   *entity.Config
	pages map[string]*entity.Page
}

func newConfig(cfg *entity.Config,
	pages map[string]*entity.Page) *BotConfig {
	return &BotConfig{
		cfg:   cfg,
		pages: pages,
	}
}

func (c *BotConfig) FindPage(name string) *entity.Page {
	return findPage(c.pages, name)
}

func findPage(pages map[string]*entity.Page, name string) *entity.Page {
	page, ok := pages[name]
	if !ok {
		return nil
	}

	return page
}

func (c *BotConfig) GetFirstPageName() string {
	return c.cfg.FirstPage.Name
}

func (c *BotConfig) findButton(page *entity.Page, name string) *entity.Button {
	for i, btn := range page.StartKeyboard.Buttons {
		if btn.Handler == name {
			return &page.StartKeyboard.Buttons[i]
		}
	}

	return nil
}

func (c *BotConfig) findHandlerPage(handlerName string) *entity.Page {
	for _, page := range c.pages {
		if page.Handler == handlerName {
			return page
		}
	}

	return nil
}

// FindHandler finds and returns handler and handler's page name
func (c *BotConfig) FindHandler(curPage *entity.Page, data entity.Data) (entity.Handler, *entity.Page, error) {
	switch data.Type {
	case entity.DataTypeMsg:
		return curPage.HandlerFn, curPage, nil
	case entity.DataTypeCmd:
		// check if it is GoTo page command
		handlerName := data.Content // GetHandlerName()
		// if err != nil {
		// 	return nil, nil, fmt.Errorf("failed find Handler: %w", err)
		// }

		p := c.findHandlerPage(handlerName)
		if p == nil {
			// try to find handler among buttons
			b := c.findButton(curPage, handlerName)
			if b == nil {
				// button not found
				// it is custom button and should be handled in page handler
				return curPage.HandlerFn, curPage, nil
			} else {
				return b.HandlerFn, curPage, nil
			}
		} else {
			return p.HandlerFn, p, nil
		}
	}

	return nil, nil, fmt.Errorf("FindHandler unexpected error")
}

func (c *BotConfig) GetWorkerHandlers() []entity.Worker {
	return c.cfg.Workers
}
