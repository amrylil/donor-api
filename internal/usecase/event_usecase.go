package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// --- Interface ---
type EventUsecase interface {
	Create(ctx context.Context, req dto.EventRequest) (entity.Event, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Event, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, req dto.EventRequest) (entity.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Implementation ---
type eventUsecaseImpl struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) EventUsecase {
	return &eventUsecaseImpl{repo: repo}
}

func (uc *eventUsecaseImpl) Create(ctx context.Context, req dto.EventRequest) (entity.Event, error) {
	var event entity.Event
	copier.Copy(&event, &req)

	event.Slug = helper.GenerateSlug(req.EventName)
	err := uc.repo.Save(ctx, &event)
	return event, err
}

func (uc *eventUsecaseImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Event, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *eventUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Event, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *eventUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.EventRequest) (entity.Event, error) {
	event, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Event{}, err
	}

	copier.Copy(&event, &req)

	return uc.repo.Update(ctx, event)
}

func (uc *eventUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
