package persistence

import (
	"context"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockRepositoryImpl struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) repository.StockRepository {
	return &stockRepositoryImpl{db: db}
}

func (r *stockRepositoryImpl) Save(ctx context.Context, stock *entity.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

func (r *stockRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.Stock, int64, error) {
	var stocks []entity.Stock
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Stock{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (r *stockRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Stock, error) {
	var stock entity.Stock
	err := r.db.WithContext(ctx).First(&stock, id).Error
	return stock, err
}

func (r *stockRepositoryImpl) Update(ctx context.Context, stock entity.Stock) (entity.Stock, error) {
	err := r.db.WithContext(ctx).Save(&stock).Error
	return stock, err
}

func (r *stockRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Stock{}, id).Error
}
