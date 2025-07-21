package handler

import (
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	userUsecase usecase.UserUsecase
}

func NewProfileHandler(userUC usecase.UserUsecase) *ProfileHandler {
	return &ProfileHandler{
		userUsecase: userUC,
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, userDetail, err := h.userUsecase.GetProfile(c, userID.(uuid.UUID))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	userResponse := dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	profileResponse := dto.ProfileResponse{
		User: userResponse,
	}

	if userDetail != nil {
		profileResponse.Details = &dto.UserDetailResponse{
			ID:            userDetail.ID.String(),
			UserID:        userDetail.UserID.String(),
			FullName:      userDetail.FullName,
			NIK:           userDetail.NIK,
			Gender:        userDetail.Gender,
			DateOfBirth:   userDetail.DateOfBirth,
			BloodType:     userDetail.BloodType,
			Rhesus:        userDetail.Rhesus,
			PhoneNumber:   userDetail.PhoneNumber,
			Address:       userDetail.Address,
			IsActiveDonor: userDetail.IsActiveDonor,
			CreatedAt:     userDetail.CreatedAt,
			UpdatedAt:     userDetail.UpdatedAt,
		}
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Profile retrieved successfully", profileResponse)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, err := h.userUsecase.UpdateProfile(c, userID.(uuid.UUID), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := dto.UserResponse{
		ID:    updatedUser.ID.String(),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  updatedUser.Role,
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Profile updated successfully", userResponse)
}

// --- User Detail Handlers ---

func (h *ProfileHandler) CreateMyDetail(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UserDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.userUsecase.CreateUserDetail(c, userID.(uuid.UUID), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccessResponse(c, http.StatusCreated, "User detail created successfully", result)
}

func (h *ProfileHandler) GetMyDetail(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	result, err := h.userUsecase.GetUserDetailByUserID(c, userID.(uuid.UUID))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "User detail not found")
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "User detail retrieved successfully", result)
}

func (h *ProfileHandler) UpdateMyDetail(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.UserDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.userUsecase.UpdateUserDetail(c, userID.(uuid.UUID), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "User detail updated successfully", result)
}
