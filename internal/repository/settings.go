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
	settingsTable = "user_settings"
)

func (r *Repo) SaveSettings(ctx context.Context, userID int64, settings []byte) error {
	userSettings := entity.UserSettings{
		UserID:   userID,
		Settings: settings,
	}

	if err := r.db.WithContext(ctx).Table(settingsTable).Save(&userSettings).Error; err != nil {
		return fmt.Errorf("SaveSettings: %w", err)
	}

	return nil
}

func (r *Repo) GetSettings(ctx context.Context, userID int64) ([]byte, error) {
	var settings entity.UserSettings

	err := r.db.WithContext(ctx).
		Table(settingsTable).
		Where("user_id = ?", userID).
		First(&settings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) ||
		errors.Is(err, logger.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return settings.Settings, nil
}
