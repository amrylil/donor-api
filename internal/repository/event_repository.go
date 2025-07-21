package repository

import (
	"context"
	"donor-api/internal/entity"

	"github.com/google/uuid"
)

type EventRepository interface {
	Save(ctx context.Context, event *entity.Event) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Event, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Event, error)
	Update(ctx context.Context, event entity.Event) (entity.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
