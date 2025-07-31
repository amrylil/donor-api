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

type EventHandler struct {
	usecase usecase.EventUsecase
}

func NewEventHandler(usecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{usecase: usecase}
}

// Create godoc
// @Summary      Create a new event
// @Description  Menambahkan acara (event) baru ke sistem
// @Tags         Events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.EventRequest    true  "Data Acara Baru"
// @Success      201   {object}  dto.SuccessWrapper  "Acara berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper    "Request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /events [post]
func (h *EventHandler) Create(c *gin.Context) {
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.EventResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusCreated, "Event created successfully", res)
}

// GetAll godoc
// @Summary      Get all events
// @Description  Mengambil daftar semua acara dengan paginasi
// @Tags         Events
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar acara"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /events [get]
func (h *EventHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.usecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var itemResponses []dto.EventResponse
	copier.Copy(&itemResponses, &items)

	// ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}

	paginatedResponse := dto.PaginatedResponse[dto.EventResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved events", paginatedResponse)
}

// GetByID godoc
// @Summary      Get event by ID
// @Description  Mengambil satu data acara berdasarkan ID
// @Tags         Events
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Acara"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil data acara"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      404  {object}  dto.ErrorWrapper    "Data tidak ditemukan"
// @Router       /events/{id} [get]
func (h *EventHandler) GetByID(c *gin.Context) {
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

	var res dto.EventResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved event", res)
}

// Update godoc
// @Summary      Update an event
// @Description  Memperbarui acara yang sudah ada berdasarkan ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string            true  "ID Acara"  format(uuid)
// @Param        body  body      dto.EventRequest  true  "Data Acara yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper  "Acara berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper    "Format ID atau request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /events/{id} [put]
func (h *EventHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.EventResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Event updated successfully", res)
}

// Delete godoc
// @Summary      Delete an event
// @Description  Menghapus acara berdasarkan ID
// @Tags         Events
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Acara"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Acara berhasil dihapus"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /events/{id} [delete]
func (h *EventHandler) Delete(c *gin.Context) {
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
	helper.SendSuccessResponse(c, http.StatusOK, "Event deleted successfully", "")
}
