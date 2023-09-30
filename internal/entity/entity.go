package entity

import (
	"encoding/json"
	"errors"
	"time"
)

type (
	Input interface {
		GetKeyboard() Keyboard
		GetUserID() int64
		GetData() Data
		GetCache() SessionData
		GetCurrentPhraseNum() int
		GetCurrentDay() int

		CreateOutput() Output
	}

	Output interface {
		SetKeyboard(kbd Keyboard) Output
		SetMessage(msg string) Output
		SetError(err error)
		SetData(data string)
		SetUserID(userID int64) Output
		SetCache(cache SessionData) Output
		SetGoToStart() Output
		SetCurrentPhraseNum(num int) Output
		SetCurrentDay(day int) Output

		GetError() error
		GetKeyboard() Keyboard
		GetMessage() string
		GetUserID() int64
		GetCache() SessionData
		GetGoToStart() bool
		GetCurrentPhraseNum() int
		GetCurrentDay() int
	}

	Handler       func(Input) Output
	WorkerHandler func(Worker) ([]Output, error)

	Config struct {
		Pages   []Page   `json:"pages"`
		Workers []Worker `json:"workers"`
		//Schedulers []Scheduler `json:"schedulers"`
		FirstPage *Page
	}

	Button struct {
		Text    string `json:"text"`
		Handler string `json:"handler"`

		ExternalPagePtr *Page
		HandlerFn       Handler
	}

	Keyboard struct {
		Buttons []Button `json:"buttons"`
	}

	Page struct {
		Name          string   `json:"name"`
		First         bool     `json:"first"`
		StartKeyboard Keyboard `json:"start_keyboard"`
		Handler       string   `json:"handler"`

		HandlerFn Handler
	}

	Worker struct {
		Name          string `json:"name"`
		Period        string `json:"period"`
		StartTime     string `json:"start_time"`
		Duration      string `json:"duration"`
		WorkerHandler string `json:"worker_handler"`
		Page          string `json:"page"`

		PageEntity *Page
		HandlerFn  WorkerHandler
	}

	// Scheduler struct {
	// 	Name    string   `json:"name"`
	// 	Handler string   `json:"handler"`
	// 	Period  Duration `json:"period"`

	// 	HandlerFn Handler
	// }

	User struct {
		ID     int `json:"id"`
		ChatID int `json:"chat_id"`
	}

	DataType uint8

	Data struct {
		Content string
		Type    DataType
	}
)

const (
	DataTypeMsg DataType = 1
	DataTypeCmd DataType = 2
)

type Duration time.Duration

var (
	ErrDataCmd = errors.New("data cmd not found")
)

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	dur, err := time.ParseDuration(v)
	if err != nil {
		return err
	}

	*d = Duration(dur)

	return nil
}
