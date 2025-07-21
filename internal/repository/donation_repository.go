package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid"
)

type DonationRepository interface {
	Save(ctx context.Context, donation *entity.Donation) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Donation, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Donation, error)
	Update(ctx context.Context, donation entity.Donation) (entity.Donation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
