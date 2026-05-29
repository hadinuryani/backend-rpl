package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

type AntrianRepository struct {
	db *sql.DB
}

func NewAntrianRepository(db *sql.DB) *AntrianRepository {
	return &AntrianRepository{db: db}
}

func (r *AntrianRepository) CountToday(ctx context.Context, dateStr string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM antrian WHERE tanggal_kunjungan = ?`, dateStr).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count today antrian: %w", err)
	}
	return count, nil
}

func (r *AntrianRepository) Create(ctx context.Context, pasienID int, tanggal time.Time, noAntrian, keluhan string) (*models.Antrian, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO antrian (pasien_id, tanggal_kunjungan, no_antrian, keluhan) VALUES (?, ?, ?, ?)`,
		pasienID, tanggal, noAntrian, keluhan,
	)
	if err != nil {
		return nil, fmt.Errorf("create antrian: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, int(id))
}

func (r *AntrianRepository) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Antrian, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM antrian WHERE pasien_id=?`, pasienID).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, pasien_id, tanggal_kunjungan, no_antrian, keluhan, status, created_at, updated_at
		 FROM antrian WHERE pasien_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		pasienID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Antrian
	for rows.Next() {
		a := models.Antrian{}
		var tgl time.Time
		rows.Scan(&a.ID, &a.PasienID, &tgl, &a.NoAntrian, &a.Keluhan, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		a.TanggalKunjungan = models.DateOnly{Time: tgl}
		list = append(list, a)
	}
	return list, total, nil
}

func (r *AntrianRepository) FindByID(ctx context.Context, id int) (*models.Antrian, error) {
	a := &models.Antrian{}
	var tgl time.Time
	var golDarah, alamat, noWa, jenisKelamin sql.NullString
	var isHamil sql.NullBool
	err := r.db.QueryRowContext(ctx,
		`SELECT a.id, a.pasien_id, a.tanggal_kunjungan, a.no_antrian, a.keluhan, a.status, a.created_at, a.updated_at,
		        p.nama_lengkap, p.golongan_darah, p.alamat, p.no_wa, p.jenis_kelamin, TIMESTAMPDIFF(YEAR, p.tanggal_lahir, CURDATE()) as umur, p.is_hamil
		 FROM antrian a JOIN pasien p ON p.id=a.pasien_id WHERE a.id=?`, id,
	).Scan(&a.ID, &a.PasienID, &tgl, &a.NoAntrian, &a.Keluhan, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		&a.NamaPasien, &golDarah, &alamat, &noWa, &jenisKelamin, &a.Umur, &isHamil)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.TanggalKunjungan = models.DateOnly{Time: tgl}
	if golDarah.Valid {
		a.GolonganDarah = golDarah.String
	}
	if alamat.Valid {
		a.Alamat = alamat.String
	}
	if noWa.Valid {
		a.NoWa = noWa.String
	}
	if jenisKelamin.Valid {
		a.JenisKelamin = jenisKelamin.String
	}
	if isHamil.Valid {
		a.IsHamil = isHamil.Bool
	}
	return a, nil
}

func (r *AntrianRepository) FindByDateAndStatus(ctx context.Context, dateStr string, status string, limit, offset int) ([]models.Antrian, int, error) {
	var total int
	if status != "" {
		r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM antrian WHERE tanggal_kunjungan=? AND status=?`, dateStr, status).Scan(&total)
	} else {
		r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM antrian WHERE tanggal_kunjungan=?`, dateStr).Scan(&total)
	}

	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT a.id, a.pasien_id, a.tanggal_kunjungan, a.no_antrian, a.keluhan, a.status, a.created_at, a.updated_at,
		         p.nama_lengkap, TIMESTAMPDIFF(YEAR, p.tanggal_lahir, CURDATE()) as umur
		         FROM antrian a JOIN pasien p ON p.id=a.pasien_id
		         WHERE a.tanggal_kunjungan=? AND a.status=? ORDER BY a.no_antrian ASC LIMIT ? OFFSET ?`
		args = []interface{}{dateStr, status, limit, offset}
	} else {
		query = `SELECT a.id, a.pasien_id, a.tanggal_kunjungan, a.no_antrian, a.keluhan, a.status, a.created_at, a.updated_at,
		         p.nama_lengkap, TIMESTAMPDIFF(YEAR, p.tanggal_lahir, CURDATE()) as umur
		         FROM antrian a JOIN pasien p ON p.id=a.pasien_id
		         WHERE a.tanggal_kunjungan=? ORDER BY a.no_antrian ASC LIMIT ? OFFSET ?`
		args = []interface{}{dateStr, limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Antrian
	for rows.Next() {
		a := models.Antrian{}
		var tgl time.Time
		rows.Scan(&a.ID, &a.PasienID, &tgl, &a.NoAntrian, &a.Keluhan, &a.Status, &a.CreatedAt, &a.UpdatedAt, &a.NamaPasien, &a.Umur)
		a.TanggalKunjungan = models.DateOnly{Time: tgl}
		list = append(list, a)
	}
	return list, total, nil
}

func (r *AntrianRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE antrian SET status=? WHERE id=?`, status, id)
	return err
}

func (r *AntrianRepository) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int, status string) error {
	_, err := tx.ExecContext(ctx, `UPDATE antrian SET status=? WHERE id=?`, status, id)
	return err
}

func (r *AntrianRepository) GetDashboardStats(ctx context.Context, dateStr string) (total, waiting, done int, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT 
			COUNT(*), 
			COALESCE(SUM(CASE WHEN status='menunggu' THEN 1 ELSE 0 END), 0), 
			COALESCE(SUM(CASE WHEN status='selesai' THEN 1 ELSE 0 END), 0)
		 FROM antrian WHERE tanggal_kunjungan=?`, dateStr,
	).Scan(&total, &waiting, &done)
	return
}

func (r *AntrianRepository) GetWeeklyVisitCounts(ctx context.Context, days int) ([]models.WeeklyVisit, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	// Inisialisasi map untuk N hari terakhir dengan jumlah 0
	counts := make(map[string]int)
	order := make([]string, days)
	for i := 0; i < days; i++ {
		d := now.AddDate(0, 0, -i)
		dateStr := d.Format("2006-01-02")
		counts[dateStr] = 0
		order[days-1-i] = dateStr // urutan kronologis (terlama ke terbaru)
	}

	startDate := order[0]
	rows, err := r.db.QueryContext(ctx,
		`SELECT tanggal_kunjungan, COUNT(*) as jumlah
		 FROM antrian
		 WHERE tanggal_kunjungan >= ?
		 GROUP BY tanggal_kunjungan`, startDate,
	)
	if err != nil {
		return nil, fmt.Errorf("get weekly visit counts query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tgl time.Time
		var count int
		if err := rows.Scan(&tgl, &count); err != nil {
			return nil, fmt.Errorf("scan weekly visit count: %w", err)
		}
		dateStr := tgl.Format("2006-01-02")
		if _, exists := counts[dateStr]; exists {
			counts[dateStr] = count
		}
	}

	hariIndo := map[time.Weekday]string{
		time.Monday:    "Sen",
		time.Tuesday:   "Sel",
		time.Wednesday: "Rab",
		time.Thursday:  "Kam",
		time.Friday:    "Jum",
		time.Saturday:  "Sab",
		time.Sunday:    "Min",
	}

	result := make([]models.WeeklyVisit, days)
	for i, dateStr := range order {
		t, _ := time.Parse("2006-01-02", dateStr)
		result[i] = models.WeeklyVisit{
			Tanggal: dateStr,
			Hari:    hariIndo[t.Weekday()],
			Jumlah:  counts[dateStr],
		}
	}

	return result, nil
}

