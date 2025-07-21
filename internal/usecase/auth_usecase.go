package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto" // Hanya di-import untuk LoginResult
	"donor-api/internal/entity"
	"donor-api/internal/infrastructure/security"
	"donor-api/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
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

// --- PERBAIKAN DI FUNGSI REGISTER ---
func (u *authUsecaseImpl) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error) {
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecaseImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !security.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	result := &dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	return result, nil
}
