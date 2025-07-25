package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/entity"
	"donor-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, *entity.UserDetail, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UserRequest) (entity.User, error)

	// user detail
	CreateUserDetail(ctx context.Context, userID uuid.UUID, req dto.UserDetailRequest) (entity.UserDetail, error)
	GetUserDetailByUserID(ctx context.Context, userID uuid.UUID) (entity.UserDetail, error)
	UpdateUserDetail(ctx context.Context, userID uuid.UUID, req dto.UserDetailRequest) (entity.UserDetail, error)
}

type userUsecaseImpl struct {
	userRepo repository.UserRepository
}

// NewAuthUsecase membuat implementasi baru untuk AuthUsecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepo: userRepo,
	}
}

func (uc *userUsecaseImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, *entity.UserDetail, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	userDetail, err := uc.userRepo.FindDetailByUserID(ctx, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil, nil
		}
		return nil, nil, err
	}

	return user, &userDetail, nil
}

func (uc *userUsecaseImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UserRequest) (entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return entity.User{}, errors.New("user not found")
	}

	user.Name = req.Name

	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		return entity.User{}, err
	}

	return *user, nil
}

func (uc *userUsecaseImpl) CreateUserDetail(ctx context.Context, userID uuid.UUID, req dto.UserDetailRequest) (entity.UserDetail, error) {
	_, err := uc.userRepo.FindDetailByUserID(ctx, userID)
	if err == nil {
		return entity.UserDetail{}, errors.New("user detail already exists")
	}

	userDetail := &entity.UserDetail{
		UserID:        userID,
		FullName:      req.FullName,
		NIK:           req.NIK,
		Gender:        req.Gender,
		DateOfBirth:   req.DateOfBirth,
		BloodType:     req.BloodType,
		Rhesus:        req.Rhesus,
		PhoneNumber:   req.PhoneNumber,
		Address:       req.Address,
		IsActiveDonor: req.IsActiveDonor,
	}

	err = uc.userRepo.SaveDetail(ctx, userDetail)
	return *userDetail, err
}

func (uc *userUsecaseImpl) GetUserDetailByUserID(ctx context.Context, userID uuid.UUID) (entity.UserDetail, error) {
	return uc.userRepo.FindDetailByUserID(ctx, userID)
}

func (uc *userUsecaseImpl) UpdateUserDetail(ctx context.Context, userID uuid.UUID, req dto.UserDetailRequest) (entity.UserDetail, error) {
	detail, err := uc.userRepo.FindDetailByUserID(ctx, userID)
	if err != nil {
		return entity.UserDetail{}, err
	}

	detail.FullName = req.FullName
	detail.NIK = req.NIK
	detail.Gender = req.Gender
	detail.DateOfBirth = req.DateOfBirth
	detail.BloodType = req.BloodType
	detail.Rhesus = req.Rhesus
	detail.PhoneNumber = req.PhoneNumber
	detail.Address = req.Address
	detail.IsActiveDonor = req.IsActiveDonor

	return uc.userRepo.UpdateDetail(ctx, detail)
}
