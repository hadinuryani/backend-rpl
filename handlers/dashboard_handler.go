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
	jadwalRepo  repositories.JadwalRepo
}

func NewDashboardHandler(ar repositories.AntrianRepo, ir repositories.InventoriRepo, pr repositories.PasienRepo, kr repositories.KlinikRepo, jr repositories.JadwalRepo) *DashboardHandler {
	return &DashboardHandler{antrianRepo: ar, invRepo: ir, pasienRepo: pr, klinikRepo: kr, jadwalRepo: jr}
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

func (h *DashboardHandler) GetChartKunjungan(c *gin.Context) {
	counts, err := h.antrianRepo.GetWeeklyVisitCounts(c.Request.Context(), 7)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data kunjungan mingguan")
		return
	}
	utils.SuccessResponse(c, "Berhasil", counts)
}

func (h *DashboardHandler) GetChartObat(c *gin.Context) {
	summary, err := h.invRepo.GetStockStatusSummary(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data stok obat")
		return
	}
	utils.SuccessResponse(c, "Berhasil", gin.H{
		"stok_aman":          summary["aman"],
		"hampir_habis":      summary["hampir_habis"],
		"stok_habis":         summary["habis"],
		"hampir_kadaluarsa": summary["hampir_kadaluarsa"],
		"kadaluarsa":        summary["kadaluarsa"],
	})
}

func (h *DashboardHandler) GetAlerts(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	todayStr := time.Now().In(loc).Format("2006-01-02")

	obatKritis, err := h.invRepo.GetCriticalMedicines(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data obat kritis")
		return
	}

	jadwalHariIni, err := h.jadwalRepo.FindTodaySchedules(c.Request.Context(), todayStr)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil jadwal hari ini")
		return
	}

	utils.SuccessResponse(c, "Berhasil", gin.H{
		"obat_kritis":     obatKritis,
		"jadwal_hari_ini": jadwalHariIni,
	})
}

