package handlers

import (
	"time"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PasienHandler struct {
	repo     *repositories.PasienRepository
	validate *validator.Validate
}

func NewPasienHandler(repo *repositories.PasienRepository) *PasienHandler {
	return &PasienHandler{repo: repo, validate: validator.New()}
}

func (h *PasienHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil || pasien == nil {
		utils.NotFound(c, "Profil pasien tidak ditemukan")
		return
	}
	utils.SuccessResponse(c, "Berhasil", pasien)
}

func (h *PasienHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	var req dto.UpdatePasienRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, "Validasi gagal", err.Error())
		return
	}

	pasien, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil || pasien == nil {
		utils.NotFound(c, "Profil pasien tidak ditemukan")
		return
	}

	nama := pasien.NamaLengkap
	if req.NamaLengkap != "" { nama = req.NamaLengkap }
	tglLahir := pasien.TanggalLahir
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			utils.BadRequest(c, "Format tanggal lahir tidak valid")
			return
		}
		tglLahir = t
	}
	jk := pasien.JenisKelamin
	if req.JenisKelamin != "" { jk = req.JenisKelamin }
	alamat := ""
	if pasien.Alamat != nil { alamat = *pasien.Alamat }
	if req.Alamat != "" { alamat = req.Alamat }
	noWa := pasien.NoWa
	if req.NoWa != "" { noWa = req.NoWa }
	var golDarah *string
	if req.GolonganDarah != "" {
		golDarah = &req.GolonganDarah
	} else {
		golDarah = pasien.GolonganDarah
	}

	err = h.repo.Update(c.Request.Context(), pasien.ID, nama, tglLahir, jk, alamat, noWa, golDarah)
	if err != nil {
		utils.InternalError(c, "Gagal memperbarui profil")
		return
	}
	utils.SuccessResponse(c, "Profil berhasil diperbarui", nil)
}
