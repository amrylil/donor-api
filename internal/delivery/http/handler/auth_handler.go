package handler

import (
	"donor-api/internal/delivery/http/dto"
	"donor-api/internal/delivery/http/helper"
	"donor-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", user)

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}
	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", res)

}
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	idtoken, err := helper.GetBearerToken(c)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid token")
		return
	}
	res, err := h.authUsecase.AuthenticateWithGoogle(c.Request.Context(), idtoken)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}
	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", res)

}
