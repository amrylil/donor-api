package persistence

import (
	"context"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type locationRepositoryImpl struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &locationRepositoryImpl{db: db}
}

func (r *locationRepositoryImpl) Save(ctx context.Context, location *entity.Location) error {
	return r.db.WithContext(ctx).Create(location).Error
}

func (r *locationRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.Location, int64, error) {
	var locations []entity.Location
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Location{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&locations).Error; err != nil {
		return nil, 0, err
	}

	return locations, total, nil
}

func (r *locationRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Location, error) {
	var location entity.Location
	err := r.db.WithContext(ctx).First(&location, id).Error
	return location, err
}
func (r *locationRepositoryImpl) FindByTenantID(ctx context.Context, limit, offset int, tenantID uuid.UUID) ([]entity.Location, int64, error) {
	var locations []entity.Location
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.Location{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).
		Where("tenant_id = ?", tenantID).
		Find(&locations).Error
	return locations, total, err
}

func (r *locationRepositoryImpl) Update(ctx context.Context, location entity.Location) (entity.Location, error) {
	err := r.db.WithContext(ctx).Save(&location).Error
	return location, err
}

func (r *locationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Location{}, id).Error
}
