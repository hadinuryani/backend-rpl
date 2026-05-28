package repositories

import (
	"context"
	"database/sql"
	"time"
)

// MonitorRepository handles database operations for visit monitoring.
type MonitorRepository struct {
	db *sql.DB
}

// NewMonitorRepository creates a new MonitorRepository instance.
func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

// GetVisitHistory retrieves paginated visit history with optional search filtering.
func (r *MonitorRepository) GetVisitHistory(ctx context.Context, from, to time.Time, search string, limit, offset int) ([]VisitRecord, int, error) {
	var total int
	r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM antrian a JOIN pasien p ON p.id=a.pasien_id
		 WHERE a.tanggal_kunjungan BETWEEN ? AND ?
		 AND (?='' OR p.nama_lengkap LIKE CONCAT('%', ?, '%') OR a.keluhan LIKE CONCAT('%', ?, '%'))`,
		from, to, search, search, search).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.pasien_id, a.no_antrian, p.nama_lengkap, a.tanggal_kunjungan, a.keluhan, a.status, r.id
		 FROM antrian a 
		 JOIN pasien p ON p.id=a.pasien_id
		 LEFT JOIN rekam_medis r ON r.antrian_id = a.id
		 WHERE a.tanggal_kunjungan BETWEEN ? AND ?
		 AND (?='' OR p.nama_lengkap LIKE CONCAT('%', ?, '%') OR a.keluhan LIKE CONCAT('%', ?, '%'))
		 ORDER BY a.tanggal_kunjungan DESC, a.no_antrian ASC
		 LIMIT ? OFFSET ?`,
		from, to, search, search, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []VisitRecord
	for rows.Next() {
		v := VisitRecord{}
		var tgl time.Time
		var rmID sql.NullInt64
		err := rows.Scan(&v.ID, &v.PasienID, &v.NoAntrian, &v.NamaPasien, &tgl, &v.Keluhan, &v.Status, &rmID)
		if err != nil {
			return nil, 0, err
		}
		v.TanggalDaftar = tgl.Format("2006-01-02")
		if rmID.Valid {
			idVal := int(rmID.Int64)
			v.RekamMedisID = &idVal
		}
		list = append(list, v)
	}
	return list, total, nil
}
