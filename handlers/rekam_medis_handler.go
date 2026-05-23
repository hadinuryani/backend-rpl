package handlers

import (
	"database/sql"
	"ic-plus-backend/dto"
	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RekamMedisHandler struct {
	service    *services.RekamMedisService
	repo       *repositories.RekamMedisRepository
	pasienRepo *repositories.PasienRepository
	db         *sql.DB
	validate   *validator.Validate
}

func NewRekamMedisHandler(svc *services.RekamMedisService, repo *repositories.RekamMedisRepository, pr *repositories.PasienRepository, db *sql.DB) *RekamMedisHandler {
	return &RekamMedisHandler{service: svc, repo: repo, pasienRepo: pr, db: db, validate: validator.New()}
}

func (h *RekamMedisHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")
	var bidanID int
	h.db.QueryRowContext(c.Request.Context(), `SELECT id FROM bidan WHERE user_id=?`, userID).Scan(&bidanID)
	if bidanID == 0 { utils.InternalError(c, "Bidan tidak ditemukan"); return }

	var clinicStatus string
	_ = h.db.QueryRowContext(c.Request.Context(), `SELECT status FROM klinik_status ORDER BY updated_at DESC LIMIT 1`).Scan(&clinicStatus)
	if clinicStatus == "tutup" {
		utils.BadRequest(c, "Klinik sedang tutup. Tidak dapat memproses pemeriksaan saat klinik tutup.")
		return
	}

	var req dto.CreateRekamMedisRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.validate.Struct(req); err != nil { utils.ValidationErrorResponse(c, "Validasi gagal", err.Error()); return }

	rmID, resepID, err := h.service.CreateWithResep(c.Request.Context(), bidanID, req)
	if err != nil { utils.InternalError(c, "Gagal: "+err.Error()); return }
	utils.CreatedResponse(c, "Rekam medis disimpan, antrian selesai", gin.H{"rekam_medis_id": rmID, "resep_id": resepID})
}

func (h *RekamMedisHandler) GetByID(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	rm, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil || rm == nil { utils.NotFound(c, "Rekam medis tidak ditemukan"); return }
	resep, _ := h.repo.FindResepByRekamMedisID(c.Request.Context(), id)
	utils.SuccessResponse(c, "Berhasil", gin.H{"rekam_medis": rm, "resep": resep})
}

func (h *RekamMedisHandler) GetMyRekamMedis(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, _ := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if pasien == nil { utils.NotFound(c, "Profil tidak ditemukan"); return }
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindByPasienID(c.Request.Context(), pasien.ID, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }

	// Enrich each record with resep data
	type RMWithResep struct {
		models.RekamMedis
		Resep []models.DetailResep `json:"resep"`
	}
	var enriched []RMWithResep
	for _, rm := range list {
		item := RMWithResep{RekamMedis: rm}
		resep, _ := h.repo.FindResepByRekamMedisID(c.Request.Context(), rm.ID)
		if resep != nil && len(resep.Details) > 0 {
			item.Resep = resep.Details
		}
		enriched = append(enriched, item)
	}

	utils.PaginatedResponse(c, "Berhasil", enriched, utils.BuildMeta(total, p))
}

func (h *RekamMedisHandler) Update(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	var req dto.UpdateRekamMedisRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	var td, kj, ct *string
	if req.TekananDarah != "" { td = &req.TekananDarah }
	if req.KondisiJanin != "" { kj = &req.KondisiJanin }
	if req.CatatanTambahan != "" { ct = &req.CatatanTambahan }
	err := h.repo.Update(c.Request.Context(), id, req.KeluhanUtama, td, req.BeratBadan, req.TinggiFundusUteri, kj, ct)
	if err != nil { utils.InternalError(c, "Gagal memperbarui"); return }
	utils.SuccessResponse(c, "Berhasil diperbarui", nil)
}

// GetRekamMedisByPasienID returns all rekam medis for a given patient (used by bidan).
func (h *RekamMedisHandler) GetRekamMedisByPasienID(c *gin.Context) {
	pasienID := parseID(c, "id")
	if pasienID == 0 { return }

	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindByPasienID(c.Request.Context(), pasienID, p.Limit, p.Offset)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data rekam medis")
		return
	}

	// Enrich each record with resep data
	type RMWithResep struct {
		models.RekamMedis
		Resep []models.DetailResep `json:"resep"`
	}
	var enriched []RMWithResep
	for _, rm := range list {
		item := RMWithResep{RekamMedis: rm}
		resep, _ := h.repo.FindResepByRekamMedisID(c.Request.Context(), rm.ID)
		if resep != nil && len(resep.Details) > 0 {
			item.Resep = resep.Details
		}
		enriched = append(enriched, item)
	}

	utils.PaginatedResponse(c, "Berhasil", enriched, utils.BuildMeta(total, p))
}
