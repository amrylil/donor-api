package handler

import (
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BloodRequestHandler struct {
	usecase usecase.BloodRequestUsecase
}

func NewBloodRequestHandler(usecase usecase.BloodRequestUsecase) *BloodRequestHandler {
	return &BloodRequestHandler{usecase: usecase}
}

// Create godoc
// @Summary      Create a new blood request
// @Description  Menambahkan permintaan darah baru ke sistem
// @Tags         Blood Requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.BloodRequestRequest  true  "Data Permintaan Darah"
// @Success      201   {object}  dto.SuccessWrapper       "Permintaan darah berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper         "Request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper         "Terjadi kesalahan internal"
// @Router       /blood-requests [post]
func (h *BloodRequestHandler) Create(c *gin.Context) {
	var req dto.BloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "BloodRequest created successfully", res)
}

// GetAll godoc
// @Summary      Get all blood requests
// @Description  Mengambil daftar semua permintaan darah dengan paginasi
// @Tags         Blood Requests
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar permintaan darah"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /blood-requests [get]
func (h *BloodRequestHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginatedResponse, err := h.usecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved blood_requests", paginatedResponse)
}

// GetByID godoc
// @Summary      Get blood request by ID
// @Description  Mengambil satu permintaan darah berdasarkan ID
// @Tags         Blood Requests
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Permintaan Darah"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil data permintaan darah"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      404  {object}  dto.ErrorWrapper    "Data tidak ditemukan"
// @Router       /blood-requests/{id} [get]
func (h *BloodRequestHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	res, err := h.usecase.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Record not found")
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved blood_request", res)
}

// Update godoc
// @Summary      Update a blood request
// @Description  Memperbarui permintaan darah yang sudah ada berdasarkan ID
// @Tags         Blood Requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                   true  "ID Permintaan Darah"  format(uuid)
// @Param        body  body      dto.BloodRequestRequest  true  "Data Permintaan Darah yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper       "Permintaan darah berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper         "Format ID atau request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper         "Terjadi kesalahan internal"
// @Router       /blood-requests/{id} [put]
func (h *BloodRequestHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req dto.BloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "BloodRequest updated successfully", res)
}

// Delete godoc
// @Summary      Delete a blood request
// @Description  Menghapus permintaan darah berdasarkan ID
// @Tags         Blood Requests
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Permintaan Darah"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Permintaan darah berhasil dihapus"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /blood-requests/{id} [delete]
func (h *BloodRequestHandler) Delete(c *gin.Context) {
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
	helper.SendSuccessResponse(c, http.StatusOK, "BloodRequest deleted successfully", "")
}
