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

type StockHandler struct {
	usecase usecase.StockUsecase
}

func NewStockHandler(usecase usecase.StockUsecase) *StockHandler {
	return &StockHandler{usecase: usecase}
}

// Create godoc
// @Summary      Create a new stock
// @Description  Menambahkan data stok darah baru ke sistem
// @Tags         Stocks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.StockRequest    true  "Data Stok Baru"
// @Success      201   {object}  dto.SuccessWrapper  "Stok berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper    "Request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /stocks [post]
func (h *StockHandler) Create(c *gin.Context) {
	var req dto.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.StockResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusCreated, "Stock created successfully", res)
}

// GetAll godoc
// @Summary      Get all stocks
// @Description  Mengambil daftar semua stok darah dengan paginasi
// @Tags         Stocks
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar stok"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /stocks [get]
func (h *StockHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.usecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var itemResponses []dto.StockResponse
	copier.Copy(&itemResponses, &items)

	// ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}

	paginatedResponse := dto.PaginatedResponse[dto.StockResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved stocks", paginatedResponse)
}

// GetByID godoc
// @Summary      Get stock by ID
// @Description  Mengambil satu data stok darah berdasarkan ID
// @Tags         Stocks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Stok"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil data stok"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      404  {object}  dto.ErrorWrapper    "Data tidak ditemukan"
// @Router       /stocks/{id} [get]
func (h *StockHandler) GetByID(c *gin.Context) {
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

	var res dto.StockResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved stock", res)
}

// Update godoc
// @Summary      Update a stock
// @Description  Memperbarui data stok darah yang sudah ada berdasarkan ID
// @Tags         Stocks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string            true  "ID Stok"  format(uuid)
// @Param        body  body      dto.StockRequest  true  "Data Stok yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper  "Stok berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper    "Format ID atau request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /stocks/{id} [put]
func (h *StockHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req dto.StockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.StockResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Stock updated successfully", res)
}

// Delete godoc
// @Summary      Delete a stock
// @Description  Menghapus data stok darah berdasarkan ID
// @Tags         Stocks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Stok"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Stok berhasil dihapus"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /stocks/{id} [delete]
func (h *StockHandler) Delete(c *gin.Context) {
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
	helper.SendSuccessResponse(c, http.StatusOK, "Stock deleted successfully", "")
}
