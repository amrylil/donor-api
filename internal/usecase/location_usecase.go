package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/entity"
	"donor-api/internal/repository"
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// --- Interface ---
type LocationUsecase interface {
	Create(ctx context.Context, req dto.LocationRequest, tenantID uuid.UUID) (*entity.Location, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Location, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Location, error)
	Update(ctx context.Context, id uuid.UUID, req dto.LocationRequest) (entity.Location, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserLocation(ctx context.Context, lat float64, lon float64) ([]dto.LocationByUserResponse, error)
}

// --- Implementation ---
type locationUsecaseImpl struct {
	repo repository.LocationRepository
}

func NewLocationUsecase(repo repository.LocationRepository) LocationUsecase {
	return &locationUsecaseImpl{repo: repo}
}

func (uc *locationUsecaseImpl) Create(ctx context.Context, req dto.LocationRequest, tenantID uuid.UUID) (*entity.Location, error) {
	location := &entity.Location{}
	copier.Copy(&location, &req)

	location.Slug = helper.GenerateSlug(req.LocationName)

	log.Print("tenantID: ", tenantID)

	location.TenantID = tenantID

	err := uc.repo.Save(ctx, location)
	if err != nil {
		return nil, err
	}
	return location, err
}

func (uc *locationUsecaseImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Location, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *locationUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Location, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *locationUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.LocationRequest) (entity.Location, error) {
	location, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Location{}, err
	}

	copier.Copy(&location, &req)

	return uc.repo.Update(ctx, location)
}

func (uc *locationUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *locationUsecaseImpl) GetAllByUserLocation(ctx context.Context, lat float64, lon float64) ([]dto.LocationByUserResponse, error) {
	locations, _, err := uc.repo.FindAll(ctx, 1000, 0) // Ambil semua lokasi
	if err != nil {
		return nil, err
	}

	var responses []dto.LocationByUserResponse
	for _, loc := range locations {
		distance := helper.Haversine(lat, lon, *loc.Latitude, *loc.Longitude)
		responses = append(responses, dto.LocationByUserResponse{
			Name:     loc.LocationName,
			Address:  loc.Address,
			Lat:      *loc.Latitude,
			Lon:      *loc.Longitude,
			Distance: distance,
		})
	}

	return responses, nil
}
