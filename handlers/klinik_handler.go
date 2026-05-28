package handlers

import (
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type KlinikHandler struct {
	repo     repositories.KlinikRepo
	bidanRepo repositories.BidanRepo
}

func NewKlinikHandler(repo repositories.KlinikRepo, bidanRepo repositories.BidanRepo) *KlinikHandler {
	return &KlinikHandler{repo: repo, bidanRepo: bidanRepo}
}

func (h *KlinikHandler) GetStatus(c *gin.Context) {
	status, catatan, err := h.repo.GetStatus(c.Request.Context())
	if err != nil {
		utils.SuccessResponse(c, "Berhasil", gin.H{"status": "tutup", "catatan": ""})
		return
	}
	utils.SuccessResponse(c, "Berhasil", gin.H{"status": status, "catatan": catatan})
}

func (h *KlinikHandler) SetStatus(c *gin.Context) {
	var req struct {
		Status  string `json:"status" binding:"required"`
		Catatan string `json:"catatan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Data tidak valid")
		return
	}
	if req.Status != "buka" && req.Status != "tutup" && req.Status != "libur" {
		utils.BadRequest(c, "Status harus 'buka', 'tutup', atau 'libur'")
		return
	}

	userID := c.GetInt("user_id")
	bidanID, err := h.bidanRepo.FindBidanIDByUserID(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, "Bidan tidak ditemukan")
		return
	}

	if err := h.repo.SetStatus(c.Request.Context(), bidanID, req.Status, req.Catatan); err != nil {
		utils.InternalError(c, "Gagal memperbarui status")
		return
	}

	utils.SuccessResponse(c, "Status klinik berhasil diperbarui", gin.H{"status": req.Status})
}
