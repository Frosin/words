package usecase

import (
	"context"
	"test/internal/entity"
	"test/internal/repository"
)

type Usecase interface {
	CreatePhrase(ctx context.Context, userID int64, phrase string) error
	UpdatePhrase(ctx context.Context, userID int64, phrase string, sentence string) error
	SetPhrasesNum(ctx context.Context, userID int64, num int) error
	GetPhrasesNum(ctx context.Context, userID int64) (int, error)
	GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error)
	GetPhraseInfo(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error)

	GetSession(ctx context.Context, userID int64) (*entity.Session, error)
	SaveSession(ctx context.Context, session entity.Session) error
}

type Uc struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &Uc{
		repo,
	}
}
