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
