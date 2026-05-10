package handlers

import (
	"database/sql"
	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type KlinikHandler struct {
	db *sql.DB
}

func NewKlinikHandler(db *sql.DB) *KlinikHandler {
	return &KlinikHandler{db: db}
}

func (h *KlinikHandler) GetStatus(c *gin.Context) {
	var status, catatan string
	var catatanNull sql.NullString
	err := h.db.QueryRowContext(c.Request.Context(),
		`SELECT status, catatan FROM klinik_status ORDER BY updated_at DESC LIMIT 1`,
	).Scan(&status, &catatanNull)

	if err != nil {
		utils.SuccessResponse(c, "Berhasil", gin.H{"status": "tutup", "catatan": ""})
		return
	}
	if catatanNull.Valid {
		catatan = catatanNull.String
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
	var bidanID int
	err := h.db.QueryRowContext(c.Request.Context(), `SELECT id FROM bidan WHERE user_id = ?`, userID).Scan(&bidanID)
	if err != nil {
		utils.InternalError(c, "Bidan tidak ditemukan")
		return
	}

	result, err := h.db.ExecContext(c.Request.Context(),
		`UPDATE klinik_status SET status = ?, catatan = ? WHERE bidan_id = ?`,
		req.Status, req.Catatan, bidanID,
	)
	if err != nil {
		utils.InternalError(c, "Gagal memperbarui status")
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_, err = h.db.ExecContext(c.Request.Context(),
			`INSERT INTO klinik_status (bidan_id, status, catatan) VALUES (?, ?, ?)`,
			bidanID, req.Status, req.Catatan,
		)
		if err != nil {
			utils.InternalError(c, "Gagal menyimpan status")
			return
		}
	}

	utils.SuccessResponse(c, "Status klinik berhasil diperbarui", gin.H{"status": req.Status})
}
