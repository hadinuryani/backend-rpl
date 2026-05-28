package handlers

import (
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type ResepHandler struct {
	rmRepo repositories.RekamMedisRepo
}

func NewResepHandler(rmRepo repositories.RekamMedisRepo) *ResepHandler {
	return &ResepHandler{rmRepo: rmRepo}
}

func (h *ResepHandler) GetByRekamMedisID(c *gin.Context) {
	rmID := parseID(c, "rekam_medis_id")
	if rmID == 0 { return }
	resep, err := h.rmRepo.FindResepByRekamMedisID(c.Request.Context(), rmID)
	if err != nil || resep == nil {
		utils.NotFound(c, "Resep tidak ditemukan")
		return
	}
	utils.SuccessResponse(c, "Berhasil", resep)
}
