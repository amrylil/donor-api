package persistence

import (
  "context"
  "donor-api/internal/entity"
  "donor-api/internal/repository"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type tenantRepositoryImpl struct {
  db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) repository.TenantRepository {
  return &tenantRepositoryImpl{db: db}
}

func (r *tenantRepositoryImpl) Save(ctx context.Context, tenant *entity.Tenant) error {
  return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.Tenant, int64, error) {
  var tenants []entity.Tenant
  var total int64

  if err := r.db.WithContext(ctx).Model(&entity.Tenant{}).Count(&total).Error; err != nil {
    return nil, 0, err
  }

  if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&tenants).Error; err != nil {
    return nil, 0, err
  }

  return tenants, total, nil
}

func (r *tenantRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Tenant, error) {
  var tenant entity.Tenant
  // GORM dapat mencari berdasarkan primary key secara langsung.
  err := r.db.WithContext(ctx).First(&tenant, id).Error
  return tenant, err
}

func (r *tenantRepositoryImpl) Update(ctx context.Context, tenant entity.Tenant) (entity.Tenant, error) {
  err := r.db.WithContext(ctx).Save(&tenant).Error
  return tenant, err
}

func (r *tenantRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
  // GORM dapat menghapus berdasarkan primary key secara langsung.
  return r.db.WithContext(ctx).Delete(&entity.Tenant{}, id).Error
}
