package handlers

import (
	"database/sql"
	"time"

	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	antrianRepo *repositories.AntrianRepository
	invRepo     *repositories.InventoriRepository
	db          *sql.DB
}

func NewDashboardHandler(ar *repositories.AntrianRepository, ir *repositories.InventoriRepository, db *sql.DB) *DashboardHandler {
	return &DashboardHandler{antrianRepo: ar, invRepo: ir, db: db}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)

	total, waiting, done, err := h.antrianRepo.GetDashboardStats(c.Request.Context(), today)
	if err != nil { utils.InternalError(c, "Gagal mengambil statistik"); return }

	critical, _ := h.invRepo.CountCritical(c.Request.Context())

	var klinikStatus string
	err = h.db.QueryRowContext(c.Request.Context(),
		`SELECT status FROM klinik_status ORDER BY updated_at DESC LIMIT 1`,
	).Scan(&klinikStatus)
	if err != nil { klinikStatus = "tutup" }

	utils.SuccessResponse(c, "Berhasil", gin.H{
		"total_pasien_hari_ini": total,
		"antrian_menunggu":      waiting,
		"antrian_selesai":       done,
		"stok_kritis":           critical,
		"klinik_status":         klinikStatus,
	})
}
