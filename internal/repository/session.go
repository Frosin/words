package repository

import (
	"context"
	"errors"
	"fmt"
	"test/internal/entity"
	"test/internal/metrics"
	"time"

	"github.com/jinzhu/gorm"
	"gorm.io/gorm/logger"
)

const (
	sessionsTable = "sessions"
)

func (r *Repo) GetSession(ctx context.Context, userID int64) (*entity.Session, error) {
	var session entity.Session

	defer func(begin time.Time) {
		delta := time.Since(begin).Seconds()
		metrics.WordsRequestDuration.WithLabelValues("get_session").Observe(delta)
	}(time.Now())

	err := r.db.WithContext(ctx).
		Table(sessionsTable).
		Where("user_id = ?", userID).
		First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, logger.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Repo) GetSessions(ctx context.Context, userIDs []int64) ([]*entity.Session, error) {
	var sessions []*entity.Session

	defer func(begin time.Time) {
		delta := time.Since(begin).Seconds()
		metrics.WordsRequestDuration.WithLabelValues("get_sessions").Observe(delta)
	}(time.Now())

	err := r.db.WithContext(ctx).
		Table(sessionsTable).
		Where("user_id IN (?)", userIDs).
		Find(&sessions).Error

	if errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, logger.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *Repo) SaveSession(ctx context.Context, session entity.Session) error {
	defer func(begin time.Time) {
		delta := time.Since(begin).Seconds()
		metrics.WordsRequestDuration.WithLabelValues("save_session").Observe(delta)
	}(time.Now())

	if err := r.db.WithContext(ctx).Table(sessionsTable).Save(&session).Error; err != nil {
		return fmt.Errorf("SaveSession: %w", err)
	}

	return nil
}
