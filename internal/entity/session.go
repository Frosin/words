package entity

import (
	"encoding/json"
	"fmt"
)

type Session struct {
	UserID           int    `gorm:"column:user_id;primaryKey" json:"user_id"`
	CurrentPage      string `gorm:"column:current_page" json:"current_page"`
	LastMsgID        int    `gorm:"column:last_msg_id" json:"last_message_id"`
	ChatID           int64  `gorm:"column:chat_id" json:"chat_id"`
	LastMsgHasKbd    bool   `gorm:"column:last_msg_has_kbd" json:"last_msg_has_kbd"`
	CurrentPhraseNum int    `gorm:"column:current_phrase_num" json:"current_phrase_num"`
	CurrentDay       int    `gorm:"column:current_day" json:"current_day"`

	Data RawDataType `gorm:"column:data" json:"data"`
}

type RawDataType string

type SessionData map[string]any

func NewSessionData() SessionData {
	return make(SessionData)
}

func (r RawDataType) GetSessionData() (SessionData, error) {
	sd := make(SessionData)

	if r == "" {
		return sd, nil
	}

	if err := json.Unmarshal([]byte(r), &sd); err != nil {
		return nil, fmt.Errorf("GetSessionData: %w", err)
	}

	return sd, nil
}

func (s SessionData) ToRaw() (RawDataType, error) {
	jsn, err := json.Marshal(s)

	return RawDataType(jsn), err
}
