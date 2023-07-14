package entity

type OutputObject struct {
	msg       string
	data      string
	kbd       Keyboard
	err       error
	userID    int64
	cache     SessionData
	goToStart bool
	// fields for session
	currentPhraseNum int
	currentDay       int
}

func NewOutput() Output {
	return &OutputObject{}
}

func (o *OutputObject) SetKeyboard(kbd Keyboard) Output {
	o.kbd = kbd

	return o
}

func (o *OutputObject) SetMessage(msg string) Output {
	o.msg = msg

	return o
}

func (o *OutputObject) SetError(err error) {
	o.err = err
}

func (o *OutputObject) SetData(data string) {
	o.data = data
}

func (o *OutputObject) SetUserID(userID int64) Output {
	o.userID = userID

	return o
}

func (o *OutputObject) GetError() error {
	return o.err
}

func (o *OutputObject) GetKeyboard() Keyboard {
	return o.kbd
}

func (o *OutputObject) GetMessage() string {
	return o.msg
}

func (o *OutputObject) GetUserID() int64 {
	return o.userID
}

func (o *OutputObject) SetCache(cache SessionData) Output {
	o.cache = cache

	return o
}

func (o *OutputObject) GetCache() SessionData {
	return o.cache
}

func (o *OutputObject) SetGoToStart() Output {
	o.goToStart = true

	return o
}

func (o *OutputObject) GetGoToStart() bool {
	return o.goToStart
}

func (o *OutputObject) SetCurrentPhraseNum(num int) Output {
	o.currentPhraseNum = num

	return o
}

func (o *OutputObject) SetCurrentDay(day int) Output {
	o.currentDay = day

	return o
}

func (o *OutputObject) GetCurrentPhraseNum() int {
	return o.currentPhraseNum
}

func (o *OutputObject) GetCurrentDay() int {
	return o.currentDay
}
