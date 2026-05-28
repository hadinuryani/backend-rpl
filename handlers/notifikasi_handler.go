package handlers

import (
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type NotifikasiHandler struct {
	repo       repositories.NotifikasiRepo
	pasienRepo repositories.PasienRepo
}

func NewNotifikasiHandler(repo repositories.NotifikasiRepo, pr repositories.PasienRepo) *NotifikasiHandler {
	return &NotifikasiHandler{repo: repo, pasienRepo: pr}
}

func (h *NotifikasiHandler) GetMyNotifikasi(c *gin.Context) {
	userID := c.GetInt("user_id")
	pasien, _ := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if pasien == nil { utils.NotFound(c, "Profil tidak ditemukan"); return }
	p := utils.GetPaginationParams(c)
	list, total, err := h.repo.FindByPasienID(c.Request.Context(), pasien.ID, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}

func (h *NotifikasiHandler) MarkAsRead(c *gin.Context) {
	id := parseID(c, "id"); if id == 0 { return }
	userID := c.GetInt("user_id")
	pasien, _ := h.pasienRepo.FindByUserID(c.Request.Context(), userID)
	if pasien == nil { utils.NotFound(c, "Profil tidak ditemukan"); return }
	if err := h.repo.MarkAsRead(c.Request.Context(), id, pasien.ID); err != nil {
		utils.InternalError(c, "Gagal memperbarui"); return
	}
	utils.SuccessResponse(c, "Notifikasi ditandai sudah dibaca", nil)
}
