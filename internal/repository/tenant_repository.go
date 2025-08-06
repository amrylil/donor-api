package repository

import (
  "context"
  "donor-api/internal/entity"

  "github.com/google/uuid"
)

type TenantRepository interface {
  Save(ctx context.Context, tenant *entity.Tenant) error
  FindAll(ctx context.Context, limit, offset int) ([]entity.Tenant, int64, error)
  FindByID(ctx context.Context, id uuid.UUID) (entity.Tenant, error)
  Update(ctx context.Context, tenant entity.Tenant) (entity.Tenant, error)
  Delete(ctx context.Context, id uuid.UUID) error
}
