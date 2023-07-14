package entity

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Phrase struct {
	gorm.Model

	Phrase string `gorm:"column:phrase"`
	LangID uint8  `gorm:"column:lang_id"`
	UserID int64  `gorm:"column:user_id"`
	Epoch  uint8  `gorm:"column:epoch"`

	Meta datatypes.JSON `gorm:"column:meta"`
}

type PhraseMeta struct {
	Sentences []string `json:"sentences"`
}

type UserSettings struct {
	UserID         int64 `gorm:"column:user_id;primaryKey"`
	PhraseDayLimit uint8 `gorm:"column:phrase_day_limit" json:"phrase_day_limit"`
}

// type Settings struct {
// 	Langs []LangSettings `json:"langs"`
// }

// type LangSettings struct {
// 	LangID uint64 `json:"lang_id"`
// }

// func (s Settings) Serialize() ([]byte, error) {
// 	return json.Marshal(s)
// }

// func DeserializeSettings(data []byte) (*Settings, error) {
// 	sets := Settings{}
// 	err := json.Unmarshal(data, &sets)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &sets, nil
// }

func (s PhraseMeta) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func DeserializePhraseMeta(data []byte) (*PhraseMeta, error) {
	meta := PhraseMeta{}
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
