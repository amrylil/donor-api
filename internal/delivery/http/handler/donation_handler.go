package handler

import (
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type DonationHandler struct {
	usecase usecase.DonationUsecase
}

func NewDonationHandler(usecase usecase.DonationUsecase) *DonationHandler {
	return &DonationHandler{usecase: usecase}
}

// Create godoc
// @Summary      Create a new donation
// @Description  Menambahkan data donasi baru ke sistem
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateDonationRequest  true  "Data Donasi Baru"
// @Success      201   {object}  dto.SuccessWrapper         "Donasi berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper           "Request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper           "Terjadi kesalahan internal"
// @Router       /donations [post]
func (h *DonationHandler) Create(c *gin.Context) {
	var req dto.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := helper.GetContextValue(c, "userID")

	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		log.Print("find error when get user_id from the context: ", err.Error())
		return
	}

	role, err := helper.GetRoleFromContext(c)

	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		log.Print("find error when get role from the context: ", err.Error())
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), req, userID, *role)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.DonationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusCreated, "Donation created successfully", res)
}

// GetAll godoc
// @Summary      Get all donations
// @Description  Mengambil daftar semua data donasi dengan paginasi
// @Tags         Donations
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar donasi"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /donations [get]
func (h *DonationHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.usecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var itemResponses []dto.DonationResponse
	copier.Copy(&itemResponses, &items)

	// ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}

	paginatedResponse := dto.PaginatedResponse[dto.DonationResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved donations", paginatedResponse)
}

// GetByID godoc
// @Summary      Get donation by ID
// @Description  Mengambil satu data donasi berdasarkan ID
// @Tags         Donations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Donasi"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil data donasi"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      404  {object}  dto.ErrorWrapper    "Data tidak ditemukan"
// @Router       /donations/{id} [get]
func (h *DonationHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	result, err := h.usecase.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Record not found")
		return
	}

	var res dto.DonationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved donation", res)
}

// Update godoc
// @Summary      Update a donation
// @Description  Memperbarui data donasi yang sudah ada berdasarkan ID
// @Tags         Donations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "ID Donasi"  format(uuid)
// @Param        body  body      dto.UpdateDonationRequest  true  "Data Donasi yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper         "Donasi berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper           "Format ID atau request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper           "Terjadi kesalahan internal"
// @Router       /donations/{id} [put]
func (h *DonationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req dto.UpdateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.DonationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Donation updated successfully", res)
}

// Delete godoc
// @Summary      Delete a donation
// @Description  Menghapus data donasi berdasarkan ID
// @Tags         Donations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Donasi"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Donasi berhasil dihapus"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /donations/{id} [delete]
func (h *DonationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = h.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Donation deleted successfully", "")
}
