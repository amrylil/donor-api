package usecase

import (
	"context"
	"donor-api/internal/delivery/http/dto" // Hanya di-import untuk LoginResult
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/entity"
	"donor-api/internal/infrastructure/security"
	"donor-api/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"

	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest, role string) (*entity.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	AuthenticateWithGoogle(ctx context.Context, idTokenString string) (*dto.LoginResponse, error)
}

type authUsecaseImpl struct {
	userRepo    repository.UserRepository
	tenantRepo  repository.TenantRepository
	jwtService  *security.JWTService
	webClientID string
}

func NewAuthUsecase(userRepo repository.UserRepository, tenantRepo repository.TenantRepository, jwtService *security.JWTService, webClientID string) AuthUsecase {
	return &authUsecaseImpl{
		userRepo:    userRepo,
		tenantRepo:  tenantRepo,
		jwtService:  jwtService,
		webClientID: webClientID,
	}
}

func (u *authUsecaseImpl) Register(ctx context.Context, req dto.RegisterRequest, role string) (*entity.User, error) {
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
	var locationID *uuid.UUID
	if req.LocationID != nil {
		id, err := uuid.Parse(*req.LocationID)
		if err != nil {
			return nil, fmt.Errorf("invalid location id: %w", err)
		}
		locationID = &id
	}

	newTenant := &entity.Tenant{
		Name: req.Name,
		Slug: helper.GenerateSlug(req.Name),
	}
	if err := u.tenantRepo.Save(ctx, newTenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	tenantID := &newTenant.ID

	user := &entity.User{
		Name:          req.Name,
		Email:         &req.Email,
		Password:      &hashedPassword,
		Role:          role,
		AccountStatus: "claimed",
		TenantID:      tenantID,
		LocationID:    locationID,
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

	if !security.CheckPasswordHash(req.Password, *user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Role, *user.TenantID)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: *user.Email,

		Role: user.Role,
	}

	result := &dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	return result, nil
}

func (u *authUsecaseImpl) AuthenticateWithGoogle(ctx context.Context, idTokenString string) (*dto.LoginResponse, error) {
	payload, err := idtoken.Validate(ctx, idTokenString, u.webClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid id token: %w", err)
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	fmt.Printf("Verifikasi Google berhasil untuk: %s (%s)\n", name, email)

	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mencari pengguna: %w", err)
		}

		fmt.Printf("Pengguna tidak ditemukan, membuat akun baru via Google...\n")

		newTenant := &entity.Tenant{
			Name: name,
			Slug: helper.GenerateSlug(name),
		}
		if err := u.tenantRepo.Save(ctx, newTenant); err != nil {
			return nil, fmt.Errorf("failed to create tenant: %w", err)
		}
		tenantID := &newTenant.ID

		newUser := &entity.User{
			Name:          name,
			Email:         &email,
			AccountStatus: "claimed",
			TenantID:      tenantID,
			Role:          "donor",
		}
		if err := u.userRepo.Save(ctx, newUser); err != nil {
			return nil, fmt.Errorf("gagal membuat user: %w", err)
		}

		user = newUser
	} else {
		fmt.Print("Pengguna ditemukan di sistem: ", user.Email)
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Role, *user.TenantID)

	if err != nil {
		return nil, err
	}

	userResponse := dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: *user.Email,
		Role:  user.Role,
	}

	result := &dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	return result, nil
}
