package handlers

import (
	"database/sql"
	"time"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BidanHandler struct {
	pasienRepo repositories.PasienRepo
	userRepo   repositories.UserRepo
	bidanRepo  repositories.BidanRepo
	db         *sql.DB
	validate   *validator.Validate
}

func NewBidanHandler(pasienRepo repositories.PasienRepo, userRepo repositories.UserRepo, bidanRepo repositories.BidanRepo, db *sql.DB) *BidanHandler {
	return &BidanHandler{pasienRepo: pasienRepo, userRepo: userRepo, bidanRepo: bidanRepo, db: db, validate: validator.New()}
}

func (h *BidanHandler) GetAllPasien(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	p := utils.GetPaginationParams(c)

	patients, total, err := h.pasienRepo.FindAll(c.Request.Context(), search, p.Limit, p.Offset)
	if err != nil {
		utils.InternalError(c, "Gagal mengambil data pasien")
		return
	}
	utils.PaginatedResponse(c, "Berhasil", patients, utils.BuildMeta(total, p))
}

func (h *BidanHandler) GetPasienByID(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 { return }

	pasien, err := h.pasienRepo.FindByID(c.Request.Context(), id)
	if err != nil || pasien == nil {
		utils.NotFound(c, "Pasien tidak ditemukan")
		return
	}
	utils.SuccessResponse(c, "Berhasil", pasien)
}

func (h *BidanHandler) CreatePasien(c *gin.Context) {
	var req dto.CreatePasienByBidanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, "Validasi gagal", err.Error())
		return
	}

	tglLahir, err := time.Parse("2006-01-02", req.TanggalLahir)
	if err != nil {
		utils.BadRequest(c, "Format tanggal lahir tidak valid")
		return
	}

	tx, err := h.db.BeginTx(c.Request.Context(), nil)
	if err != nil {
		utils.InternalError(c, "Gagal memulai transaksi")
		return
	}
	defer tx.Rollback()

	dummyEmail := req.NoWa + "@ic-plus.local"
	dummyHash := "$2a$12$placeholder"

	user, err := h.userRepo.CreateTx(c.Request.Context(), tx, dummyEmail, dummyHash, "pasien")
	if err != nil {
		utils.InternalError(c, "Gagal membuat akun pasien")
		return
	}

	var golDarah *string
	if req.GolonganDarah != "" { golDarah = &req.GolonganDarah }

	pasien, err := h.pasienRepo.CreateByBidan(c.Request.Context(), tx, req.NamaLengkap, tglLahir, req.JenisKelamin, req.Alamat, req.NoWa, golDarah, user.ID)
	if err != nil {
		utils.InternalError(c, "Gagal menyimpan data pasien")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.InternalError(c, "Gagal menyimpan")
		return
	}

	utils.CreatedResponse(c, "Pasien berhasil ditambahkan", pasien)
}

func (h *BidanHandler) UpdatePasien(c *gin.Context) {
	id := parseID(c, "id")
	if id == 0 { return }

	var req dto.UpdatePasienRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}

	pasien, err := h.pasienRepo.FindByID(c.Request.Context(), id)
	if err != nil || pasien == nil {
		utils.NotFound(c, "Pasien tidak ditemukan")
		return
	}

	nama := pasien.NamaLengkap
	if req.NamaLengkap != "" { nama = req.NamaLengkap }
	tglLahir := pasien.TanggalLahir
	if req.TanggalLahir != "" {
		t, _ := time.Parse("2006-01-02", req.TanggalLahir)
		tglLahir = t
	}
	jk := pasien.JenisKelamin
	if req.JenisKelamin != "" { jk = req.JenisKelamin }
	alamat := ""
	if pasien.Alamat != nil { alamat = *pasien.Alamat }
	if req.Alamat != "" { alamat = req.Alamat }
	noWa := pasien.NoWa
	if req.NoWa != "" { noWa = req.NoWa }
	golDarah := pasien.GolonganDarah
	if req.GolonganDarah != "" { golDarah = &req.GolonganDarah }

	err = h.pasienRepo.Update(c.Request.Context(), id, nama, tglLahir, jk, alamat, noWa, golDarah)
	if err != nil {
		utils.InternalError(c, "Gagal memperbarui pasien")
		return
	}
	utils.SuccessResponse(c, "Pasien berhasil diperbarui", nil)
}
