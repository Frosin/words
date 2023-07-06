package usecase

import (
	"context"
	"fmt"
	"test/internal/entity"
)

const (
	DefaultPhrasesNum = 3
)

func (u *Uc) SetPhrasesNum(ctx context.Context, userID int64, num int) error {
	settings := entity.Settings{
		Langs: []entity.LangSettings{
			{
				LangID:    0,
				PhraseNum: uint8(num),
			},
		},
	}

	serialized, err := settings.Serialize()
	if err != nil {
		return err
	}

	return u.repo.SaveSettings(ctx, userID, serialized)
}

func (u *Uc) GetPhrasesNum(ctx context.Context, userID int64) (int, error) {
	serialized, err := u.repo.GetSettings(ctx, userID)
	if err != nil {
		return 0, err
	}

	// if settings not found, we should create it
	if serialized == nil && err == nil {
		err := u.SetPhrasesNum(ctx, userID, DefaultPhrasesNum)
		if err != nil {
			return 0, err
		}

		return DefaultPhrasesNum, nil
	}

	settings, err := entity.DeserializeSettings(serialized)
	if err != nil {
		return 0, err
	}

	if len(settings.Langs) == 0 {
		return 0, fmt.Errorf("not found langs settings of user=%d", userID)
	}

	return int(settings.Langs[0].PhraseNum), nil
}
