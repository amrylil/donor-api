package persistence

import (
  "context"
  "donor-api/internal/entity"
  "donor-api/internal/repository"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type bloodRequestRepositoryImpl struct {
  db *gorm.DB
}

func NewBloodRequestRepository(db *gorm.DB) repository.BloodRequestRepository {
  return &bloodRequestRepositoryImpl{db: db}
}

func (r *bloodRequestRepositoryImpl) Save(ctx context.Context, bloodRequest *entity.BloodRequest) error {
  return r.db.WithContext(ctx).Create(bloodRequest).Error
}

func (r *bloodRequestRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.BloodRequest, int64, error) {
  var bloodRequests []entity.BloodRequest
  var total int64

  if err := r.db.WithContext(ctx).Model(&entity.BloodRequest{}).Count(&total).Error; err != nil {
    return nil, 0, err
  }

  if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&bloodRequests).Error; err != nil {
    return nil, 0, err
  }

  return bloodRequests, total, nil
}

func (r *bloodRequestRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.BloodRequest, error) {
  var bloodRequest entity.BloodRequest
  // GORM dapat mencari berdasarkan primary key secara langsung.
  err := r.db.WithContext(ctx).First(&bloodRequest, id).Error
  return bloodRequest, err
}

func (r *bloodRequestRepositoryImpl) Update(ctx context.Context, bloodRequest entity.BloodRequest) (entity.BloodRequest, error) {
  err := r.db.WithContext(ctx).Save(&bloodRequest).Error
  return bloodRequest, err
}

func (r *bloodRequestRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
  // GORM dapat menghapus berdasarkan primary key secara langsung.
  return r.db.WithContext(ctx).Delete(&entity.BloodRequest{}, id).Error
}
