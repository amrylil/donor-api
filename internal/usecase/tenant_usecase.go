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
type TenantUsecase interface {
	Create(ctx context.Context, req dto.TenantRequest) (dto.TenantResponse, error)
	FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.TenantResponse], error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.TenantResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.TenantRequest) (dto.TenantResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Implementation ---
type tenantUsecaseImpl struct {
	repo repository.TenantRepository
}

func NewTenantUsecase(repo repository.TenantRepository) TenantUsecase {
	return &tenantUsecaseImpl{repo: repo}
}

func (uc *tenantUsecaseImpl) Create(ctx context.Context, req dto.TenantRequest) (dto.TenantResponse, error) {
	var tenant entity.Tenant
	var res dto.TenantResponse

	copier.Copy(&tenant, &req)

	if err := uc.repo.Save(ctx, &tenant); err != nil {
		return res, err
	}

	copier.Copy(&res, &tenant)
	res.ID = tenant.ID.String()
	return res, nil
}

func (uc *tenantUsecaseImpl) FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.TenantResponse], error) {
	offset := (page - 1) * limit
	var paginatedResponse dto.PaginatedResponse[dto.TenantResponse]

	items, total, err := uc.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return paginatedResponse, err
	}

	var itemResponses []dto.TenantResponse
	copier.Copy(&itemResponses, &items)

	// ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}

	paginatedResponse = dto.PaginatedResponse[dto.TenantResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	return paginatedResponse, nil
}

func (uc *tenantUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.TenantResponse, error) {
	var res dto.TenantResponse
	tenant, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return res, err
	}

	copier.Copy(&res, &tenant)
	res.ID = tenant.ID.String()
	return res, nil
}

func (uc *tenantUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.TenantRequest) (dto.TenantResponse, error) {
	var res dto.TenantResponse
	tenant, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return res, err
	}

	copier.Copy(&tenant, &req)

	updatedTenant, err := uc.repo.Update(ctx, tenant)
	if err != nil {
		return res, err
	}

	copier.Copy(&res, &updatedTenant)
	res.ID = updatedTenant.ID.String()
	return res, nil
}

func (uc *tenantUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
