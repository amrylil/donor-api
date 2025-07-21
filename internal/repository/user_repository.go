package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error

	// user detail
	SaveDetail(ctx context.Context, userDetail *entity.UserDetail) error
	FindDetailByUserID(ctx context.Context, userID uuid.UUID) (entity.UserDetail, error)
	UpdateDetail(ctx context.Context, userDetail entity.UserDetail) (entity.UserDetail, error)
}
