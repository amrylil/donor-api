package handler

import (
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type ProfileHandler struct {
	userUsecase usecase.UserUsecase
}

func NewProfileHandler(userUC usecase.UserUsecase) *ProfileHandler {
	return &ProfileHandler{
		userUsecase: userUC,
	}
}

// GetAll godoc
// @Summary      Get all users
// @Description  Mengambil daftar semua lokasi dengan paginasi
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar lokasi"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /users [get]
func (h *ProfileHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.userUsecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var itemResponses []dto.UserResponse
	copier.Copy(&itemResponses, &items)

	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}
	paginatedResponse := dto.PaginatedResponse[dto.UserResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved users", paginatedResponse)
}

// GetProfile godoc
// @Summary      Get current user's profile
// @Description  Mengambil profil dasar dan detail dari pengguna yang sedang login
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.SuccessWrapper  "Profil berhasil diambil"
// @Failure      401  {object}  dto.ErrorWrapper    "Tidak terautentikasi"
// @Failure      404  {object}  dto.ErrorWrapper    "Profil tidak ditemukan"
// @Router       /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, err := helper.GetContextValue(c, "userID")
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, detail, err := h.userUsecase.GetProfile(c.Request.Context(), *userID)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response := dto.ProfileResponse{
		User:    *user,
		Details: detail,
	}

	helper.SendSuccessResponse(c, http.StatusOK, "User profile retrieved successfully", response)
}

// UpdateProfile godoc
// @Summary      Update current user's profile
// @Description  Memperbarui informasi dasar (nama, email) dari pengguna yang sedang login
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserRequest     true  "Data Profil yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper  "Profil berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper    "Request tidak valid"
// @Failure      401   {object}  dto.ErrorWrapper    "Tidak terautentikasi"
// @Failure      500   {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /profile [put]
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

// CreateMyDetail godoc
// @Summary      Create my user detail
// @Description  Membuat profil detail (NIK, alamat, dll.) untuk pengguna yang sedang login. Hanya bisa dibuat sekali.
// @Tags         Profile Details
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserDetailRequest  true  "Data Detail Profil"
// @Success      201   {object}  dto.SuccessWrapper     "Detail profil berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper       "Request tidak valid"
// @Failure      401   {object}  dto.ErrorWrapper       "Tidak terautentikasi"
// @Failure      500   {object}  dto.ErrorWrapper       "Terjadi kesalahan internal"
// @Router       /profile/details [post]
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

// GetMyDetail godoc
// @Summary      Get my user detail
// @Description  Mengambil profil detail dari pengguna yang sedang login
// @Tags         Profile Details
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.SuccessWrapper  "Detail profil berhasil diambil"
// @Failure      401  {object}  dto.ErrorWrapper    "Tidak terautentikasi"
// @Failure      404  {object}  dto.ErrorWrapper    "Detail profil tidak ditemukan"
// @Router       /profile/details [get]
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

// UpdateMyDetail godoc
// @Summary      Update my user detail
// @Description  Memperbarui profil detail dari pengguna yang sedang login
// @Tags         Profile Details
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UserDetailRequest  true  "Data Detail Profil yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper     "Detail profil berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper       "Request tidak valid"
// @Failure      401   {object}  dto.ErrorWrapper       "Tidak terautentikasi"
// @Failure      500   {object}  dto.ErrorWrapper       "Terjadi kesalahan internal"
// @Router       /profile/details [put]
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
