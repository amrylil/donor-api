package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// --- Interface ---
type StockUsecase interface {
	Create(ctx context.Context, req dto.StockRequest) (entity.Stock, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Stock, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Stock, error)
	Update(ctx context.Context, id uuid.UUID, req dto.StockRequest) (entity.Stock, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Implementation ---
type stockUsecaseImpl struct {
	repo repository.StockRepository
}

func NewStockUsecase(repo repository.StockRepository) StockUsecase {
	return &stockUsecaseImpl{repo: repo}
}

func (uc *stockUsecaseImpl) Create(ctx context.Context, req dto.StockRequest) (entity.Stock, error) {
	var stock entity.Stock
	copier.Copy(&stock, &req)

	err := uc.repo.Save(ctx, &stock)
	return stock, err
}

func (uc *stockUsecaseImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Stock, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *stockUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Stock, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *stockUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.StockRequest) (entity.Stock, error) {
	stock, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Stock{}, err
	}

	copier.Copy(&stock, &req)

	return uc.repo.Update(ctx, stock)
}

func (uc *stockUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
