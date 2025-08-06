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

type TenantHandler struct {
  usecase usecase.TenantUsecase
}

func NewTenantHandler(usecase usecase.TenantUsecase) *TenantHandler {
  return &TenantHandler{usecase: usecase}
}

func (h *TenantHandler) Create(c *gin.Context) {
  var req dto.TenantRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }

  res, err := h.usecase.Create(c.Request.Context(), req)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusCreated, "Tenant created successfully", res)
}

func (h *TenantHandler) GetAll(c *gin.Context) {
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

  paginatedResponse, err := h.usecase.FindAll(c.Request.Context(), page, limit)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved tenants", paginatedResponse)
}

func (h *TenantHandler) GetByID(c *gin.Context) {
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

  helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved tenant", res)
}

func (h *TenantHandler) Update(c *gin.Context) {
  id, err := uuid.Parse(c.Param("id"))
  if err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
    return
  }
  var req dto.TenantRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }

  res, err := h.usecase.Update(c.Request.Context(), id, req)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusOK, "Tenant updated successfully", res)
}

func (h *TenantHandler) Delete(c *gin.Context) {
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
  helper.SendSuccessResponse(c, http.StatusOK, "Tenant deleted successfully", "")
}
