package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid"
)

type StockRepository interface {
	Save(ctx context.Context, stock *entity.Stock) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Stock, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Stock, error)
	Update(ctx context.Context, stock entity.Stock) (entity.Stock, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
