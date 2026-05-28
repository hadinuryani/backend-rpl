package handlers

import (
	"time"

	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	antrianRepo repositories.AntrianRepo
	invRepo     repositories.InventoriRepo
	pasienRepo  repositories.PasienRepo
	klinikRepo  repositories.KlinikRepo
}

func NewDashboardHandler(ar repositories.AntrianRepo, ir repositories.InventoriRepo, pr repositories.PasienRepo, kr repositories.KlinikRepo) *DashboardHandler {
	return &DashboardHandler{antrianRepo: ar, invRepo: ir, pasienRepo: pr, klinikRepo: kr}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	todayStr := time.Now().In(loc).Format("2006-01-02")

	total, waiting, done, err := h.antrianRepo.GetDashboardStats(c.Request.Context(), todayStr)
	if err != nil { utils.InternalError(c, "Gagal mengambil statistik"); return }

	critical, _ := h.invRepo.CountCritical(c.Request.Context())

	_, totalPatients, _ := h.pasienRepo.FindAll(c.Request.Context(), "", 1, 0)

	klinikStatus := h.klinikRepo.GetStatusString(c.Request.Context())

	utils.SuccessResponse(c, "Berhasil", gin.H{
		"total_pasien_hari_ini":  total,
		"antrian_menunggu":       waiting,
		"antrian_selesai":        done,
		"stok_obat_kritis":       critical,
		"klinik_status":          klinikStatus,
		"total_pasien_terdaftar": totalPatients,
	})
}
