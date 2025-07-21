package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid" // <-- Tambahkan import
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) // <-- Ubah tipe id
}
