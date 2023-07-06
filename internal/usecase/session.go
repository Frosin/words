package usecase

import (
	"context"
	"test/internal/entity"
)

func (u *Uc) GetSession(ctx context.Context, userID int64) (*entity.Session, error) {
	return u.repo.GetSession(ctx, userID)
}
func (u *Uc) SaveSession(ctx context.Context, session entity.Session) error {
	return u.repo.SaveSession(ctx, session)
}
