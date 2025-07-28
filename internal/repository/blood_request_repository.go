package repository

import (
  "context"
  "donor-api/internal/entity"

  "github.com/google/uuid"
)

type BloodRequestRepository interface {
  Save(ctx context.Context, bloodRequest *entity.BloodRequest) error
  FindAll(ctx context.Context, limit, offset int) ([]entity.BloodRequest, int64, error)
  FindByID(ctx context.Context, id uuid.UUID) (entity.BloodRequest, error)
  Update(ctx context.Context, bloodRequest entity.BloodRequest) (entity.BloodRequest, error)
  Delete(ctx context.Context, id uuid.UUID) error
}
