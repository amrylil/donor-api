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

type DonationHandler struct {
	usecase usecase.DonationUsecase
}

func NewDonationHandler(usecase usecase.DonationUsecase) *DonationHandler {
	return &DonationHandler{usecase: usecase}
}

func (h *DonationHandler) Create(c *gin.Context) {
	var req dto.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var res dto.DonationResponse
	copier.Copy(&res, &result)
	res.ID = result.ID.String()

	helper.SendSuccessResponse(c, http.StatusCreated, "Donation created successfully", res)
}

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
