package repository

import (
	"context"
	"errors"
	"fmt"
	"test/internal/entity"

	"github.com/jinzhu/gorm"
	"gorm.io/gorm/logger"
)

const (
	phrasesTable = "phrases"
)

func (r *Repo) GetPhrase(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error) {
	obj := entity.Phrase{}

	err := r.db.WithContext(ctx).
		Table(phrasesTable).
		Where("phrase = ?", phrase).
		Where("user_id = ?", userID).
		First(&obj).Error

	if errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, logger.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (r *Repo) SavePhrase(ctx context.Context, phrase entity.Phrase) (uint, error) {
	if err := r.db.WithContext(ctx).Table(phrasesTable).Save(&phrase).Error; err != nil {
		return 0, fmt.Errorf("SavePhrase: %w", err)
	}

	return phrase.ID, nil
}

func (r *Repo) GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error) {
	phrases := []*entity.Phrase{}

	// нужно поправить запрос
	// - ограничение на количество присылаемых сообщений в день (можно проверить уже присланные по дате апдейта)
	// - установить интервал посылки
	// - + переменная - количество прысылаемых за раз слов

	// добавить метрики на ошибки для нормальной отладки

	err := r.db.Debug().WithContext(ctx).
		Table(phrasesTable).
		//Where(`(epoch = 0 AND (julianday('now') - julianday(updated_at) ) * 24 > 2)`). // after 2 hours
		// for tests
		Where(`(epoch = 0 AND (julianday('now') - julianday(updated_at) ) > 0.0002)`).
		// Or(`(epoch = 1 AND (julianday('now') - julianday(updated_at)) * 24 > 24)`). // after 1 day
		// Or(`(epoch = 2 AND (julianday('now') - julianday(updated_at)) > 14)`).      // after 2 weeks
		// Or(`(epoch = 3 AND (julianday('now') - julianday(updated_at)) > 60)`).      // after 2 months
		//Where(`(epoch = 0 AND (julianday('now') - julianday(updated_at) ) * 24 > 2)`).
		Or(`(epoch = 1 AND (julianday('now') - julianday(updated_at)) * 24 > 24)`).
		Or(`(epoch = 2 AND (julianday('now') - julianday(updated_at)) * 24 > 24*2)`).
		Or(`(epoch = 3 AND (julianday('now') - julianday(updated_at)) * 24 > 24*3)`).
		Scan(&phrases).Error

	if err != nil {
		return nil, fmt.Errorf("GetReminderPhrases: %w", err)
	}

	return phrases, nil
}
