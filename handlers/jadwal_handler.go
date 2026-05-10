package handlers

import (
	"database/sql"
	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type JadwalHandler struct {
	service    *services.JadwalService
	repo       *repositories.JadwalRepository
	pasienRepo *repositories.PasienRepository
	db         *sql.DB
	validate   *validator.Validate
}

func NewJadwalHandler(svc *services.JadwalService, repo *repositories.JadwalRepository, pr *repositories.PasienRepository, db *sql.DB) *JadwalHandler {
	return &JadwalHandler{service: svc, repo: repo, pasienRepo: pr, db: db, validate: validator.New()}
}

func (h *JadwalHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")
	var bidanID int
	h.db.QueryRowContext(c.Request.Context(), `SELECT id FROM bidan WHERE user_id=?`, userID).Scan(&bidanID)
	if bidanID == 0 { utils.InternalError(c, "Bidan tidak ditemukan"); return }

	var req dto.CreateJadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.validate.Struct(req); err != nil { utils.ValidationErrorResponse(c, "Validasi gagal", err.Error()); return }

	jk, err := h.service.Create(c.Request.Context(), req.PasienID, bidanID, req.TanggalKontrol, req.Catatan)
	if err != nil { utils.BadRequest(c, err.Error()); return }
	utils.CreatedResponse(c, "Jadwal kontrol berhasil disimpan", jk)
}

func (h *JadwalHandler) GetAll(c *gin.Context) {
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindAll(c.Request.Context(), p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

func (h *JadwalHandler) Update(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	var req dto.UpdateJadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.service.Update(c.Request.Context(), id, req.TanggalKontrol, req.Catatan); err != nil {
		utils.InternalError(c, err.Error()); return
	}
	utils.SuccessResponse(c, "Jadwal berhasil diperbarui", nil)
}

func (h *JadwalHandler) Delete(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	if err := h.service.Delete(c.Request.Context(), id); err != nil { utils.InternalError(c, "Gagal menghapus"); return }
	utils.SuccessResponse(c, "Jadwal berhasil dihapus", nil)
}

func (h *JadwalHandler) GetMyJadwal(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, _ := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if pasien == nil { utils.NotFound(c, "Profil tidak ditemukan"); return }
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindByPasienID(c.Request.Context(), pasien.ID, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}
