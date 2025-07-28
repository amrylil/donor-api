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
type BloodRequestUsecase interface {
  Create(ctx context.Context, req dto.BloodRequestRequest) (dto.BloodRequestResponse, error)
  FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.BloodRequestResponse], error)
  FindByID(ctx context.Context, id uuid.UUID) (dto.BloodRequestResponse, error)
  Update(ctx context.Context, id uuid.UUID, req dto.BloodRequestRequest) (dto.BloodRequestResponse, error)
  Delete(ctx context.Context, id uuid.UUID) error
}

// --- Implementation ---
type bloodRequestUsecaseImpl struct {
  repo repository.BloodRequestRepository
}

func NewBloodRequestUsecase(repo repository.BloodRequestRepository) BloodRequestUsecase {
  return &bloodRequestUsecaseImpl{repo: repo}
}

func (uc *bloodRequestUsecaseImpl) Create(ctx context.Context, req dto.BloodRequestRequest) (dto.BloodRequestResponse, error) {
  var bloodRequest entity.BloodRequest
  var res dto.BloodRequestResponse

  copier.Copy(&bloodRequest, &req)

  if err := uc.repo.Save(ctx, &bloodRequest); err != nil {
    return res, err
  }

  // Salin field yang cocok, lalu atur ID secara manual.
  copier.Copy(&res, &bloodRequest)
  res.ID = bloodRequest.ID.String()
  return res, nil
}

func (uc *bloodRequestUsecaseImpl) FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.BloodRequestResponse], error) {
  offset := (page - 1) * limit
  var paginatedResponse dto.PaginatedResponse[dto.BloodRequestResponse]

  items, total, err := uc.repo.FindAll(ctx, limit, offset)
  if err != nil {
    return paginatedResponse, err
  }

  var itemResponses []dto.BloodRequestResponse
  copier.Copy(&itemResponses, &items)

  // ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
  for i := range items {
    itemResponses[i].ID = items[i].ID.String()
  }

  paginatedResponse = dto.PaginatedResponse[dto.BloodRequestResponse]{
    Data:       itemResponses,
    TotalItems: total,
    Page:       page,
    Limit:      limit,
  }
  return paginatedResponse, nil
}

func (uc *bloodRequestUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.BloodRequestResponse, error) {
  var res dto.BloodRequestResponse
  bloodRequest, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return res, err
  }

  copier.Copy(&res, &bloodRequest)
  res.ID = bloodRequest.ID.String()
  return res, nil
}

func (uc *bloodRequestUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.BloodRequestRequest) (dto.BloodRequestResponse, error) {
  var res dto.BloodRequestResponse
  bloodRequest, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return res, err
  }

  copier.Copy(&bloodRequest, &req)

  updatedBloodRequest, err := uc.repo.Update(ctx, bloodRequest)
  if err != nil {
    return res, err
  }

  copier.Copy(&res, &updatedBloodRequest)
  res.ID = updatedBloodRequest.ID.String()
  return res, nil
}

func (uc *bloodRequestUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
  _, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return err
  }
  return uc.repo.Delete(ctx, id)
}
