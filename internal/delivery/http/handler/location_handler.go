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

type LocationHandler struct {
	usecase usecase.LocationUsecase
}

func NewLocationHandler(usecase usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: usecase}
}

// Create godoc
// @Summary      Create a new location
// @Description  Menambahkan data lokasi baru ke sistem
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.LocationRequest  true  "Data Lokasi Baru"
// @Success      201   {object}  dto.SuccessWrapper   "Lokasi berhasil dibuat"
// @Failure      400   {object}  dto.ErrorWrapper     "Request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper     "Terjadi kesalahan internal"
// @Router       /locations [post]
func (h *LocationHandler) Create(c *gin.Context) {
	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.LocationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusCreated, "Location created successfully", res)
}

// GetAll godoc
// @Summary      Get all locations
// @Description  Mengambil daftar semua lokasi dengan paginasi
// @Tags         Locations
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Nomor halaman"  default(1)
// @Param        limit  query     int  false  "Jumlah item per halaman"  default(10)
// @Success      200    {object}  dto.SuccessWrapper  "Berhasil mengambil daftar lokasi"
// @Failure      500    {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /locations [get]
func (h *LocationHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.usecase.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var itemResponses []dto.LocationResponse
	copier.Copy(&itemResponses, &items)

	// ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
	for i := range items {
		itemResponses[i].ID = items[i].ID.String()
	}

	paginatedResponse := dto.PaginatedResponse[dto.LocationResponse]{
		Data:       itemResponses,
		TotalItems: total,
		Page:       page,
		Limit:      limit,
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved locations", paginatedResponse)
}

// GetByID godoc
// @Summary      Get location by ID
// @Description  Mengambil satu data lokasi berdasarkan ID
// @Tags         Locations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Lokasi"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil data lokasi"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      404  {object}  dto.ErrorWrapper    "Data tidak ditemukan"
// @Router       /locations/{id} [get]
func (h *LocationHandler) GetByID(c *gin.Context) {
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

	var res dto.LocationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved location", res)
}

// Update godoc
// @Summary      Update a location
// @Description  Memperbarui data lokasi yang sudah ada berdasarkan ID
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string               true  "ID Lokasi"  format(uuid)
// @Param        body  body      dto.LocationRequest  true  "Data Lokasi yang Diperbarui"
// @Success      200   {object}  dto.SuccessWrapper   "Lokasi berhasil diperbarui"
// @Failure      400   {object}  dto.ErrorWrapper     "Format ID atau request tidak valid"
// @Failure      500   {object}  dto.ErrorWrapper     "Terjadi kesalahan internal"
// @Router       /locations/{id} [put]
func (h *LocationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.LocationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusOK, "Location updated successfully", res)
}

// Delete godoc
// @Summary      Delete a location
// @Description  Menghapus data lokasi berdasarkan ID
// @Tags         Locations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID Lokasi"  format(uuid)
// @Success      200  {object}  dto.SuccessWrapper  "Lokasi berhasil dihapus"
// @Failure      400  {object}  dto.ErrorWrapper    "Format ID tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /locations/{id} [delete]
func (h *LocationHandler) Delete(c *gin.Context) {
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
	helper.SendSuccessResponse(c, http.StatusOK, "Location deleted successfully", "")
}

// GetAllByUserLocation godoc
// @Summary      Get nearby locations
// @Description  Mengambil daftar lokasi terurut dari yang terdekat berdasarkan latitude dan longitude pengguna
// @Tags         Locations
// @Produce      json
// @Security     BearerAuth
// @Param        lat  query     number  true  "Latitude Pengguna"
// @Param        lon  query     number  true  "Longitude Pengguna"
// @Success      200  {object}  dto.SuccessWrapper  "Berhasil mengambil daftar lokasi terdekat"
// @Failure      400  {object}  dto.ErrorWrapper    "Format latitude atau longitude tidak valid"
// @Failure      500  {object}  dto.ErrorWrapper    "Terjadi kesalahan internal"
// @Router       /locations/nearby [get]
func (h *LocationHandler) GetAllByUserLocation(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	userLat, errLat := strconv.ParseFloat(latStr, 64)
	userLon, errLon := strconv.ParseFloat(lonStr, 64)

	if errLat != nil || errLon != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid latitude or longitude format")
		return
	}

	sortedLocations, err := h.usecase.GetAllByUserLocation(c.Request.Context(), userLat, userLon)

	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved locations", sortedLocations)

}
