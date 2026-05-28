package handlers

import (
	"time"

	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	repo repositories.MonitorRepo
}

func NewMonitorHandler(repo repositories.MonitorRepo) *MonitorHandler {
	return &MonitorHandler{repo: repo}
}

func (h *MonitorHandler) GetVisitHistory(c *gin.Context) {
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")
	search := c.DefaultQuery("search", "")
	p := utils.GetPaginationParams(c)

	now := time.Now()
	from := now.AddDate(0, -1, 0)
	to := now
	if fromStr != "" { from, _ = time.Parse("2006-01-02", fromStr) }
	if toStr != "" { to, _ = time.Parse("2006-01-02", toStr) }

	list, total, err := h.repo.GetVisitHistory(c.Request.Context(), from, to, search, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}
