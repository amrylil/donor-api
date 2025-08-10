package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid"
)

type LocationRepository interface {
	Save(ctx context.Context, location *entity.Location) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Location, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Location, error)
	FindByTenantID(ctx context.Context, limit, offset int, tenantID uuid.UUID) ([]entity.Location, int64, error)
	Update(ctx context.Context, location entity.Location) (entity.Location, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
