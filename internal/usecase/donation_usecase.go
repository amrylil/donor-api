package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type DonationUsecase interface {
	Create(ctx context.Context, req dto.CreateDonationRequest, userID *uuid.UUID, role string) (entity.Donation, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Donation, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Donation, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateDonationRequest) (entity.Donation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type donationUsecaseImpl struct {
	repo repository.DonationRepository
}

func NewDonationUsecase(repo repository.DonationRepository) DonationUsecase {
	return &donationUsecaseImpl{repo: repo}
}

func (uc *donationUsecaseImpl) Create(ctx context.Context, req dto.CreateDonationRequest, userID *uuid.UUID, role string) (entity.Donation, error) {
	var donation entity.Donation
	copier.Copy(&donation, &req)

	if role != "admin" {
		donation.UserID = nil
	}

	donation.UserID = (*uuid.UUID)(userID)

	err := uc.repo.Save(ctx, &donation)
	return donation, err
}

func (uc *donationUsecaseImpl) FindAll(ctx context.Context, page, limit int) ([]entity.Donation, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.FindAll(ctx, limit, offset)
}

func (uc *donationUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Donation, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *donationUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateDonationRequest) (entity.Donation, error) {
	donation, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Donation{}, err
	}

	copier.Copy(&donation, &req)

	return uc.repo.Update(ctx, donation)
}

func (uc *donationUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
