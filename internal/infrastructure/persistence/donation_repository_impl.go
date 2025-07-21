package persistence

import (
	"context"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type donationRepositoryImpl struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) repository.DonationRepository {
	return &donationRepositoryImpl{db: db}
}

func (r *donationRepositoryImpl) Save(ctx context.Context, donation *entity.Donation) error {
	return r.db.WithContext(ctx).Create(donation).Error
}

func (r *donationRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.Donation, int64, error) {
	var donations []entity.Donation
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Donation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&donations).Error; err != nil {
		return nil, 0, err
	}

	return donations, total, nil
}

func (r *donationRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Donation, error) {
	var donation entity.Donation
	err := r.db.WithContext(ctx).First(&donation, id).Error
	return donation, err
}

func (r *donationRepositoryImpl) Update(ctx context.Context, donation entity.Donation) (entity.Donation, error) {
	err := r.db.WithContext(ctx).Save(&donation).Error
	return donation, err
}

func (r *donationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Donation{}, id).Error
}
