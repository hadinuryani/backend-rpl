package handlers

import (
	"ic-plus-backend/dto"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  *services.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service, validate: validator.New()}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, "Validasi gagal", err.Error())
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.NamaLengkap, req.TanggalLahir, req.JenisKelamin, req.Alamat, req.NoWa, req.GolonganDarah)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Registrasi berhasil", gin.H{"id": user.ID, "email": user.Email, "role": user.Role})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, "Validasi gagal", err.Error())
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "Login berhasil", dto.LoginResponse{Token: token, User: gin.H{"id": user.ID, "email": user.Email, "role": user.Role}})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless — client should discard the token
	utils.SuccessResponse(c, "Logout berhasil", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, profile, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.SuccessResponse(c, "Berhasil", dto.UserProfileResponse{ID: user.ID, Email: user.Email, Role: user.Role, Profile: profile})
}
