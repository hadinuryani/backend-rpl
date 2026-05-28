package handlers

import (
	"fmt"
	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/services"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type JadwalHandler struct {
	service        *services.JadwalService
	repo           repositories.JadwalRepo
	pasienRepo     repositories.PasienRepo
	bidanRepo      repositories.BidanRepo
	pengaturanRepo repositories.PengaturanRepo
	validate       *validator.Validate
	scheduler      *services.SchedulerService
}

func NewJadwalHandler(svc *services.JadwalService, repo repositories.JadwalRepo, pr repositories.PasienRepo, bidanRepo repositories.BidanRepo, pengaturanRepo repositories.PengaturanRepo, scheduler *services.SchedulerService) *JadwalHandler {
	return &JadwalHandler{service: svc, repo: repo, pasienRepo: pr, bidanRepo: bidanRepo, pengaturanRepo: pengaturanRepo, validate: validator.New(), scheduler: scheduler}
}

func (h *JadwalHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")
	bidanID, err := h.bidanRepo.FindBidanIDByUserID(c.Request.Context(), userID)
	if err != nil { utils.InternalError(c, "Bidan tidak ditemukan"); return }

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
	settings, err := h.pengaturanRepo.GetAll(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Gagal mengambil pengaturan")
		return
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
		if err := h.pengaturanRepo.Upsert(c.Request.Context(), "waktu_pengingat", timeStr); err != nil {
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
		if err := h.pengaturanRepo.Upsert(c.Request.Context(), "nama_klinik", req.NamaKlinik); err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan nama klinik")
			return
		}
	}

	// Update alamat_klinik if provided
	if req.AlamatKlinik != "" {
		if err := h.pengaturanRepo.Upsert(c.Request.Context(), "alamat_klinik", req.AlamatKlinik); err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan alamat klinik")
			return
		}
	}

	// Update jam_kontrol if provided
	if req.JamKontrol != "" {
		if err := h.pengaturanRepo.Upsert(c.Request.Context(), "jam_kontrol", req.JamKontrol); err != nil {
			utils.InternalError(c, "Gagal menyimpan pengaturan jam kontrol")
			return
		}
	}

	utils.SuccessResponse(c, "Pengaturan berhasil diperbarui", nil)
}
