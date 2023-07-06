package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"test/internal/entity"
	"test/internal/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository interface {
	GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error)
	GetPhrase(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error)
	SavePhrase(ctx context.Context, phrase entity.Phrase) (uint, error)

	SaveSettings(ctx context.Context, userID int64, settings []byte) error
	GetSettings(ctx context.Context, userID int64) ([]byte, error)

	GetSession(ctx context.Context, userID int64) (*entity.Session, error)
	GetSessions(ctx context.Context, userIDs []int64) ([]*entity.Session, error)
	SaveSession(ctx context.Context, session entity.Session) error
}

type Repo struct {
	db *gorm.DB
}

func NewRepository(sc service.ServiceConfig) Repository {
	var notExist bool
	fileName := sc.GetDBFileName()

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		notExist = true
	}

	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if notExist {
		db.AutoMigrate(&entity.Phrase{}, &entity.UserSettings{}, &entity.Session{})
		if err != nil {
			log.Fatal(err)
		}
	}

	repo := &Repo{
		db: db,
	}

	return repo
}
