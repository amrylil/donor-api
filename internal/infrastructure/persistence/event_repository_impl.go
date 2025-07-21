package persistence

import (
	"context"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) Save(ctx context.Context, event *entity.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.Event, int64, error) {
	var events []entity.Event
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Event, error) {
	var event entity.Event
	err := r.db.WithContext(ctx).First(&event, id).Error
	return event, err
}

func (r *eventRepositoryImpl) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	err := r.db.WithContext(ctx).Save(&event).Error
	return event, err
}

func (r *eventRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Event{}, id).Error
}
