package usecase

import (
	"context"
	"donor-api/internal/entity"
	"errors"

	"donor-api/internal/infrastructure/security"
	"donor-api/internal/repository"

	"gorm.io/gorm"
)

// AuthUsecase mendefinisikan "kontrak" untuk logika bisnis autentikasi
type AuthUsecase interface {
	Register(ctx context.Context, name, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authUsecaseImpl struct {
	userRepo   repository.UserRepository
	jwtService *security.JWTService
}

// NewAuthUsecase membuat implementasi baru untuk AuthUsecase
func NewAuthUsecase(userRepo repository.UserRepository, jwtService *security.JWTService) AuthUsecase {
	return &authUsecaseImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *authUsecaseImpl) Register(ctx context.Context, name, email, password string) (*entity.User, error) {
	_, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}
	user.Password = "" // Jangan kirim password hash ke response
	return user, nil
}

func (u *authUsecaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !security.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return u.jwtService.GenerateToken(user.ID)
}
