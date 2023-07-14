package usecase

import (
	"context"
	"test/internal/entity"
)

const (
	DefaultPhrasesNum = 3
)

func (u *Uc) SetPhrasesDayLimit(ctx context.Context, userID int64, limit int) error {
	settings := entity.UserSettings{
		UserID:         userID,
		PhraseDayLimit: uint8(limit),
	}

	return u.repo.SaveSettings(ctx, settings)
}

func (u *Uc) GetPhrasesDayLimit(ctx context.Context, userID int64) (int, error) {
	settings, err := u.repo.GetSettings(ctx, userID)
	if err != nil {
		return 0, err
	}

	// if settings not found, we should create it
	if settings == nil && err == nil {
		err := u.SetPhrasesDayLimit(ctx, userID, DefaultPhrasesNum)
		if err != nil {
			return 0, err
		}

		settings = &entity.UserSettings{
			UserID:         userID,
			PhraseDayLimit: DefaultPhrasesNum,
		}

		return DefaultPhrasesNum, nil
	}

	return int(settings.PhraseDayLimit), nil
}
