package handlers

import (
	"database/sql"
	"time"

	"ic-plus-backend/utils"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	db *sql.DB
}

func NewMonitorHandler(db *sql.DB) *MonitorHandler {
	return &MonitorHandler{db: db}
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

	var total int
	h.db.QueryRowContext(c.Request.Context(),
		`SELECT COUNT(*) FROM antrian a JOIN pasien p ON p.id=a.pasien_id
		 WHERE a.tanggal_kunjungan BETWEEN ? AND ?
		 AND (?='' OR p.nama_lengkap LIKE CONCAT('%', ?, '%') OR a.keluhan LIKE CONCAT('%', ?, '%'))`,
		from, to, search, search, search).Scan(&total)

	rows, err := h.db.QueryContext(c.Request.Context(),
		`SELECT a.id, a.pasien_id, a.no_antrian, p.nama_lengkap, a.tanggal_kunjungan, a.keluhan, a.status
		 FROM antrian a JOIN pasien p ON p.id=a.pasien_id
		 WHERE a.tanggal_kunjungan BETWEEN ? AND ?
		 AND (?='' OR p.nama_lengkap LIKE CONCAT('%', ?, '%') OR a.keluhan LIKE CONCAT('%', ?, '%'))
		 ORDER BY a.tanggal_kunjungan DESC, a.no_antrian ASC
		 LIMIT ? OFFSET ?`,
		from, to, search, search, search, p.Limit, p.Offset)
	if err != nil { utils.InternalError(c, "Gagal mengambil data"); return }
	defer rows.Close()

	type Visit struct {
		ID        int    `json:"id"`
		PasienID  int    `json:"pasien_id"`
		NoAntrian string `json:"no_antrian"`
		Nama      string `json:"nama_pasien"`
		Tanggal   string `json:"tanggal_daftar"`
		Keluhan   string `json:"keluhan"`
		Status    string `json:"status"`
	}
	var list []Visit
	for rows.Next() {
		v := Visit{}
		var tgl time.Time
		err := rows.Scan(&v.ID, &v.PasienID, &v.NoAntrian, &v.Nama, &tgl, &v.Keluhan, &v.Status)
		if err != nil {
			utils.InternalError(c, "Gagal memproses data")
			return
		}
		v.Tanggal = tgl.Format("2006-01-02")
		list = append(list, v)
	}
	utils.PaginatedResponse(c, "Berhasil", list, utils.BuildMeta(total, p))
}
