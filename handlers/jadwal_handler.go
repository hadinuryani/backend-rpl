package handlers

import (
	"database/sql"
	"fmt"
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
	scheduler  *services.SchedulerService
}

func NewJadwalHandler(svc *services.JadwalService, repo *repositories.JadwalRepository, pr *repositories.PasienRepository, db *sql.DB, scheduler *services.SchedulerService) *JadwalHandler {
	return &JadwalHandler{service: svc, repo: repo, pasienRepo: pr, db: db, validate: validator.New(), scheduler: scheduler}
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

func (h *JadwalHandler) GetWaktuPengingat(c *gin.Context) {
	rows, err := h.db.QueryContext(c.Request.Context(), `SELECT kunci, nilai FROM pengaturan`)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil pengaturan")
		return
	}
	defer rows.Close()

	settings := gin.H{
		"waktu_pengingat": "08:00",
		"nama_klinik":      "Klinik Indah Care Plus (IC+)",
		"alamat_klinik":    "Jl. Indah Care No. 45, Jakarta",
		"jam_kontrol":      "08:00 - selesai",
	}

	for rows.Next() {
		var kunci, nilai string
		if err := rows.Scan(&kunci, &nilai); err == nil {
			settings[kunci] = nilai
		}
	}
	utils.SuccessResponse(c, "Berhasil", settings)
}

func (h *JadwalHandler) UpdateWaktuPengingat(c *gin.Context) {
	var req struct {
		Waktu        string `json:"waktu_pengingat"`
		NamaKlinik   string `json:"nama_klinik"`
		AlamatKlinik string `json:"alamat_klinik"`
		JamKontrol   string `json:"jam_kontrol"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Format data tidak valid")
		return
	}

	// Update waktu_pengingat if provided
	if req.Waktu != "" {
		// Validate format HH:MM
		var hour, minute int
		_, err := fmt.Sscanf(req.Waktu, "%d:%d", &hour, &minute)
		if err != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
			utils.BadRequest(c, "Format waktu tidak valid (HH:MM)")
			return
		}
		timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
		_, err = h.db.ExecContext(c.Request.Context(),
			`INSERT INTO pengaturan (kunci, nilai) VALUES ('waktu_pengingat', ?) ON DUPLICATE KEY UPDATE nilai = ?`,
			timeStr, timeStr,
		)
		if err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan waktu pengingat")
			return
		}
		if err := h.scheduler.Reschedule(timeStr); err != nil {
			utils.InternalError(c, "Gagal menerapkan jadwal pengingat baru")
			return
		}
	}

	// Update nama_klinik if provided
	if req.NamaKlinik != "" {
		_, err := h.db.ExecContext(c.Request.Context(),
			`INSERT INTO pengaturan (kunci, nilai) VALUES ('nama_klinik', ?) ON DUPLICATE KEY UPDATE nilai = ?`,
			req.NamaKlinik, req.NamaKlinik,
		)
		if err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan nama klinik")
			return
		}
	}

	// Update alamat_klinik if provided
	if req.AlamatKlinik != "" {
		_, err := h.db.ExecContext(c.Request.Context(),
			`INSERT INTO pengaturan (kunci, nilai) VALUES ('alamat_klinik', ?) ON DUPLICATE KEY UPDATE nilai = ?`,
			req.AlamatKlinik, req.AlamatKlinik,
		)
		if err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan alamat klinik")
			return
		}
	}

	// Update jam_kontrol if provided
	if req.JamKontrol != "" {
		_, err := h.db.ExecContext(c.Request.Context(),
			`INSERT INTO pengaturan (kunci, nilai) VALUES ('jam_kontrol', ?) ON DUPLICATE KEY UPDATE nilai = ?`,
			req.JamKontrol, req.JamKontrol,
		)
		if err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan jam kontrol")
			return
		}
	}

	utils.SuccessResponse(c, "Pengaturan berhasil diperbarui", nil)
}
