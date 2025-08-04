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

// Register godoc
// @Summary      Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Data Registrasi"
// @Success      201   {object}  dto.SuccessWrapper   "User berhasil dibuat, data user ada di field 'data'"
// @Failure      400   {object}  dto.ErrorWrapper     "Request tidak valid"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), req, "user")
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", user)

}

// Register godoc
// @Summary      Register a new admin user
// @Description  Register a new user with admin role
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Data Registrasi"
// @Success      201   {object}  dto.SuccessWrapper   "Admin berhasil dibuat, data user ada di field 'data'"
// @Failure      400   {object}  dto.ErrorWrapper     "Request tidak valid"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), req, "admin")
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", user)

}

// Register godoc
// @Summary      Register a new super admin user
// @Description  Register a new user with super admin role
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Data Registrasi"
// @Success      201   {object}  dto.SuccessWrapper   "Super admin berhasil dibuat, data user ada di field 'data'"
// @Failure      400   {object}  dto.ErrorWrapper     "Request tidak valid"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterSuperAdmin(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), req, "superadmin")
	if err != nil {
		helper.SendErrorResponse(c, http.StatusNonAuthoritativeInfo, err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "User created successfully", user)
}

// Login godoc
// @Summary      User Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest    true  "Kredensial Login"
// @Success      200   {object}  dto.SuccessWrapper  "Login berhasil, token ada di field 'data'"
// @Failure      401   {object}  dto.ErrorWrapper    "Kredensial tidak valid"
// @Router       /auth/login [post]
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

// GoogleAuth godoc
// @Summary      Google Authentication
// @Description  Autentikasi pengguna menggunakan Google ID Token dan mengembalikan token login
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  dto.SuccessWrapper  "Autentikasi berhasil, token ada di field 'data'"
// @Failure      400   {object}  dto.ErrorWrapper    "Token tidak valid atau hilang"
// @Failure      401   {object}  dto.ErrorWrapper    "Autentikasi Google gagal"
// @Router       /auth/google [post]
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
