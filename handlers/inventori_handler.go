package handlers

import (
	"database/sql"
	"time"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InventoriHandler struct {
	service  *services.InventoriService
	repo     *repositories.InventoriRepository
	db       *sql.DB
	validate *validator.Validate
}

func NewInventoriHandler(svc *services.InventoriService, repo *repositories.InventoriRepository, db *sql.DB) *InventoriHandler {
	return &InventoriHandler{service: svc, repo: repo, db: db, validate: validator.New()}
}

func (h *InventoriHandler) GetAllObat(c *gin.Context) {
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindAllObat(c.Request.Context(), p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

func (h *InventoriHandler) CreateObat(c *gin.Context) {
	var req dto.CreateObatWithInventoriRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.validate.Struct(req); err != nil { utils.ValidationErrorResponse(c, "Validasi gagal", err.Error()); return }

	stokMin := req.BatasStokKritis
	if stokMin == 0 { stokMin = req.StokMinimum }
	if stokMin == 0 { stokMin = 10 }

	var tglKadaluarsa *time.Time
	if req.TanggalKadaluarsa != "" {
		dateStr := req.TanggalKadaluarsa
		if len(dateStr) > 10 {
			dateStr = dateStr[:10]
		}
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil { utils.BadRequest(c, "Format tanggal kadaluarsa tidak valid"); return }
		tglKadaluarsa = &t
	}
	var batch *string
	if req.BatchNumber != "" { batch = &req.BatchNumber }

	obat, err := h.repo.CreateObatWithInventori(c.Request.Context(), req.NamaObat, req.Kategori, req.Satuan, stokMin, req.JumlahStok, tglKadaluarsa, batch)
	if err != nil { utils.InternalError(c, "Gagal menyimpan: "+err.Error()); return }
	utils.CreatedResponse(c, "Obat berhasil ditambahkan", obat)
}

func (h *InventoriHandler) UpdateObat(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	var req dto.UpdateObatRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }

	stokMin := req.BatasStokKritis
	if stokMin == nil { stokMin = req.StokMinimum }

	var tglKadaluarsa *time.Time
	if req.TanggalKadaluarsa != "" {
		dateStr := req.TanggalKadaluarsa
		if len(dateStr) > 10 {
			dateStr = dateStr[:10]
		}
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil { utils.BadRequest(c, "Format tanggal kadaluarsa tidak valid"); return }
		tglKadaluarsa = &t
	}

	if err := h.repo.UpdateObat(c.Request.Context(), id, req.NamaObat, req.Kategori, req.Satuan, stokMin, req.JumlahStok, tglKadaluarsa); err != nil {
		utils.InternalError(c, "Gagal memperbarui"); return
	}
	utils.SuccessResponse(c, "Obat berhasil diperbarui", nil)
}

func (h *InventoriHandler) DeleteObat(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	if err := h.repo.DeleteObat(c.Request.Context(), id); err != nil { utils.InternalError(c, "Gagal menghapus"); return }
	utils.SuccessResponse(c, "Obat berhasil dihapus", nil)
}

func (h *InventoriHandler) GetInventori(c *gin.Context) {
	statusFilter := c.DefaultQuery("status", "")
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindAllInventori(c.Request.Context(), statusFilter, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

func (h *InventoriHandler) StokMasuk(c *gin.Context) {
	userID := c.GetInt("user_id")
	var bidanID int
	h.db.QueryRowContext(c.Request.Context(), `SELECT id FROM bidan WHERE user_id=?`, userID).Scan(&bidanID)
	if bidanID == 0 { utils.InternalError(c, "Bidan tidak ditemukan"); return }

	var req dto.StokMasukRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.validate.Struct(req); err != nil { utils.ValidationErrorResponse(c, "Validasi gagal", err.Error()); return }

	if err := h.service.StokMasuk(c.Request.Context(), req.InventoriID, bidanID, req.Jumlah, req.Keterangan, req.TanggalKadaluarsa, req.BatchNumber); err != nil {
		utils.BadRequest(c, err.Error()); return
	}
	utils.SuccessResponse(c, "Stok masuk berhasil dicatat", nil)
}

func (h *InventoriHandler) StokKeluar(c *gin.Context) {
	userID := c.GetInt("user_id")
	var bidanID int
	h.db.QueryRowContext(c.Request.Context(), `SELECT id FROM bidan WHERE user_id=?`, userID).Scan(&bidanID)
	if bidanID == 0 { utils.InternalError(c, "Bidan tidak ditemukan"); return }

	var req dto.StokKeluarRequest
	if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, "Data tidak valid"); return }
	if err := h.validate.Struct(req); err != nil { utils.ValidationErrorResponse(c, "Validasi gagal", err.Error()); return }

	if err := h.service.StokKeluar(c.Request.Context(), req.InventoriID, bidanID, req.Jumlah, req.Keterangan); err != nil {
		utils.BadRequest(c, err.Error()); return
	}
	utils.SuccessResponse(c, "Stok keluar berhasil dicatat", nil)
}

func (h *InventoriHandler) GetRiwayat(c *gin.Context) {
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.GetRiwayat(c.Request.Context(), p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}
