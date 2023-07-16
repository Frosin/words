package usecase

import (
	"context"
	"test/internal/backup"
	"test/internal/entity"
	"test/internal/repository"
)

type Usecase interface {
	CreatePhrase(ctx context.Context, userID int64, phrase string) error
	UpdatePhrase(ctx context.Context, userID int64, phrase string, sentence string) error
	SetPhrasesDayLimit(ctx context.Context, userID int64, limit int) error
	GetPhrasesDayLimit(ctx context.Context, userID int64) (int, error)
	GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error)
	GetPhraseInfo(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error)
	DeletePhrase(ctx context.Context, userID int64, phrase string) error
	ScheduleBackUp()

	GetSession(ctx context.Context, userID int64) (*entity.Session, error)
	SaveSession(ctx context.Context, session entity.Session) error
	ClearUsersDayLimits(ctx context.Context) error
}

type Uc struct {
	repo   repository.Repository
	dumper *backup.Dumper
}

func NewUsecase(repo repository.Repository, dumper *backup.Dumper) Usecase {
	return &Uc{
		repo,
		dumper,
	}
}

func (u *Uc) ScheduleBackUp() {
	u.dumper.ScheduleUpdate()
}
