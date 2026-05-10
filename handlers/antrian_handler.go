package handlers

import (
	"strconv"
	"time"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AntrianHandler struct {
	service     *services.AntrianService
	repo        *repositories.AntrianRepository
	pasienRepo  *repositories.PasienRepository
	validate    *validator.Validate
}

func NewAntrianHandler(service *services.AntrianService, repo *repositories.AntrianRepository, pasienRepo *repositories.PasienRepository) *AntrianHandler {
	return &AntrianHandler{service: service, repo: repo, pasienRepo: pasienRepo, validate: validator.New()}
}

// CreateAntrian handles patient visit registration.
func (h *AntrianHandler) CreateAntrian(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, err := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil || pasien == nil {
		utils.NotFound(c, "Profil pasien tidak ditemukan")
		return
	}

	var req dto.CreateAntrianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, "Validasi gagal", err.Error())
		return
	}

	antrian, err := h.service.CreateAntrian(c.Request.Context(), pasien.ID, req.TanggalKunjungan, req.Keluhan)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Pendaftaran berhasil", antrian)
}

// GetMyAntrian returns the logged-in patient's antrian history.
func (h *AntrianHandler) GetMyAntrian(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, _ := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if pasien == nil {
		utils.NotFound(c, "Profil pasien tidak ditemukan")
		return
	}
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindByPasienID(c.Request.Context(), pasien.ID, p.Limit, p.Offset)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data antrian")
		return
	}
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

// GetAntrianDetail returns a single antrian detail.
func (h *AntrianHandler) GetAntrianDetail(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 { return }

	antrian, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil || antrian == nil {
		utils.NotFound(c, "Antrian tidak ditemukan")
		return
	}
	utils.SuccessResponse(c, "Berhasil", antrian)
}

// GetTodayAntrian returns today's queue list for bidan.
func (h *AntrianHandler) GetTodayAntrian(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	status := c.DefaultQuery("status", "")
	p := utils.GetPaginationParams(c)

	list, total, err := h.repo.FindTodayByStatus(c.Request.Context(), today, status, p.Limit, p.Offset)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data antrian")
		return
	}
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

// CancelAntrian cancels an antrian (bidan).
func (h *AntrianHandler) CancelAntrian(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 { return }

	err := h.repo.UpdateStatus(c.Request.Context(), id, "batal")
	if err != nil {
		utils.InternalError(c, "Gagal membatalkan antrian")
		return
	}
	utils.SuccessResponse(c, "Antrian berhasil dibatalkan", nil)
}

// parseID extracts and validates an integer ID from URL params.
func parseID(c *gin.Context, param string) int {
	idStr := c.Param(param)
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.BadRequest(c, "ID tidak valid")
		return 0
	}
	return id
}
