package persistence

import (
	"context"
	"donor-api/internal/entity"
	"donor-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) { // <-- Ubah tipe id
	var user entity.User
	// GORM secara cerdas menangani query dengan UUID
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepositoryImpl) SaveDetail(ctx context.Context, userDetail *entity.UserDetail) error {
	return r.db.WithContext(ctx).Create(userDetail).Error
}

func (r *userRepositoryImpl) FindDetailByUserID(ctx context.Context, userID uuid.UUID) (entity.UserDetail, error) {
	var detail entity.UserDetail
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&detail).Error
	return detail, err
}

func (r *userRepositoryImpl) UpdateDetail(ctx context.Context, userDetail entity.UserDetail) (entity.UserDetail, error) {
	err := r.db.WithContext(ctx).Save(&userDetail).Error
	return userDetail, err
}
