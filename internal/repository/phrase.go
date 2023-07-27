package repository

import (
	"context"
	"errors"
	"fmt"
	"test/internal/entity"
	"time"

	"github.com/jinzhu/gorm"
	"gorm.io/gorm/logger"
)

const (
	phrasesTable = "phrases"
)

func (r *Repo) GetPhrase(ctx context.Context, userID int64, phrase string) (*entity.Phrase, error) {
	obj := entity.Phrase{}

	defer sendMetric(time.Now(), "get_phrase")

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
	defer sendMetric(time.Now(), "save_phrases")

	if err := r.db.WithContext(ctx).Table(phrasesTable).Save(&phrase).Error; err != nil {
		return 0, fmt.Errorf("SavePhrase: %w", err)
	}

	return phrase.ID, nil
}

func (r *Repo) DeletePhrase(ctx context.Context, phrase *entity.Phrase) error {
	defer sendMetric(time.Now(), "delete_phrase")

	if err := r.db.WithContext(ctx).Table(phrasesTable).Delete(&phrase).Error; err != nil {
		return fmt.Errorf("DeletePhrase: %w", err)
	}

	return nil
}

// select * from phrases p
// join sessions s on s.user_id = p.user_id
// join user_settings us on us.user_id = p.user_id
// Where
// (epoch = 0 AND (julianday('now') - julianday(updated_at) ) > 0.0002
// Or (epoch = 1 AND (julianday('now') - julianday(updated_at)) * 24 > 24)
// Or (epoch = 2 AND (julianday('now') - julianday(updated_at)) * 24 > 24*2)
// Or (epoch = 3 AND (julianday('now') - julianday(updated_at)) * 24 > 24*3))
// And
// current_word_num < phrase_day_limit
// group by p.user_id
func (r *Repo) GetReminderPhrases(ctx context.Context) ([]*entity.Phrase, error) {
	phrases := []*entity.Phrase{}

	defer sendMetric(time.Now(), "get_phrases")

	err := r.db.Debug().WithContext(ctx).
		Table(phrasesTable).
		Joins("join sessions s on s.user_id = phrases.user_id").
		Joins("join user_settings us on us.user_id = phrases.user_id").
		// for tests
		Where(
			r.db.
				Where(`(epoch = 0 AND (julianday('now') - julianday(updated_at) ) * 24 > 2)`). // after 2 hours
				Or(`(epoch = 1 AND (julianday('now') - julianday(updated_at)) * 24 > 24)`).    // after 1 day
				Or(`(epoch = 2 AND (julianday('now') - julianday(updated_at)) > 14)`).         // after 2 weeks
				Or(`(epoch = 3 AND (julianday('now') - julianday(updated_at)) > 60)`)).        // after 2 months
		// r.db.Where(`(epoch = 0 AND (julianday('now') - julianday(updated_at) ) > 0.0002)`).
		// 	Or(`(epoch = 1 AND (julianday('now') - julianday(updated_at)) > 0.0002)`)).
		//
		//Or(`(epoch = 1 AND (julianday('now') - julianday(updated_at)) * 24 > 24)`).
		//Or(`(epoch = 2 AND (julianday('now') - julianday(updated_at)) * 24 > 24*2)`).
		//Or(`(epoch = 3 AND (julianday('now') - julianday(updated_at)) * 24 > 24*3)`)).
		Where("current_phrase_num < phrase_day_limit").
		Where("phrases.deleted_at IS NULL").
		Order("epoch desc").
		Order("created_at desc").
		Scan(&phrases).Error

	if err != nil {
		return nil, fmt.Errorf("GetReminderPhrases: %w", err)
	}

	return phrases, nil
}

func (r *Repo) ClearUsersDayLimits(ctx context.Context) error {
	defer sendMetric(time.Now(), "clear_limits")

	err := r.db.Debug().WithContext(ctx).
		Exec(`
		update sessions 
		set current_day = strftime('%d','now'), current_phrase_num = 0
		where current_day != strftime('%d','now') 
		`).Error

	if err != nil {
		return fmt.Errorf("ClearUsersDayLimits: %w", err)
	}

	return nil
}
