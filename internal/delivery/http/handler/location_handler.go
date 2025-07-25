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
