package entity

type InputObject struct {
	data   Data
	kbd    Keyboard
	userID int64
	cache  SessionData
	// fields for session
	currentPhraseNum int
	currentDay       int
}

func NewInput(data Data, kbd Keyboard, userID int64, cache SessionData, currentPhraseNum, currentDay int) Input {
	return &InputObject{
		data:             data,
		kbd:              kbd,
		userID:           userID,
		cache:            cache,
		currentPhraseNum: currentPhraseNum,
		currentDay:       currentDay,
	}
}

func (i *InputObject) GetKeyboard() Keyboard {
	return i.kbd
}

func (i *InputObject) GetUserID() int64 {
	return i.userID
}

func (i *InputObject) GetData() Data {
	return i.data
}

func (i *InputObject) CreateOutput() Output {
	return NewOutput()
}

func (i *InputObject) GetCache() SessionData {
	return i.cache
}

func (i *InputObject) GetCurrentPhraseNum() int {
	return i.currentPhraseNum
}

func (i *InputObject) GetCurrentDay() int {
	return i.currentDay
}
