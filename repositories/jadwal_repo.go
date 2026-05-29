package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

type JadwalRepository struct {
	db *sql.DB
}

func NewJadwalRepository(db *sql.DB) *JadwalRepository {
	return &JadwalRepository{db: db}
}

func (r *JadwalRepository) Create(ctx context.Context, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO jadwal_kontrol (pasien_id, bidan_id, tanggal_kontrol, catatan)
		 VALUES (?,?,?,?)`,
		pasienID, bidanID, tanggal, catatan,
	)
	if err != nil {
		return nil, fmt.Errorf("create jadwal: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	jk := &models.JadwalKontrol{}
	var tgl time.Time
	err = r.db.QueryRowContext(ctx,
		`SELECT id, pasien_id, bidan_id, tanggal_kontrol, catatan, status_notifikasi, created_at, updated_at
		 FROM jadwal_kontrol WHERE id=?`, id,
	).Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt)
	if err != nil {
		return nil, err
	}
	jk.TanggalKontrol = tgl
	return jk, nil
}

func (r *JadwalRepository) CreateTx(ctx context.Context, tx *sql.Tx, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error) {
	result, err := tx.ExecContext(ctx,
		`INSERT INTO jadwal_kontrol (pasien_id, bidan_id, tanggal_kontrol, catatan)
		 VALUES (?,?,?,?)`,
		pasienID, bidanID, tanggal, catatan,
	)
	if err != nil {
		return nil, fmt.Errorf("create jadwal inside tx: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	jk := &models.JadwalKontrol{}
	var tgl time.Time
	err = tx.QueryRowContext(ctx,
		`SELECT id, pasien_id, bidan_id, tanggal_kontrol, catatan, status_notifikasi, created_at, updated_at
		 FROM jadwal_kontrol WHERE id=?`, id,
	).Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt)
	if err != nil {
		return nil, err
	}
	jk.TanggalKontrol = tgl
	return jk, nil
}

func (r *JadwalRepository) FindAll(ctx context.Context, limit, offset int) ([]models.JadwalKontrol, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jadwal_kontrol`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT jk.id, jk.pasien_id, jk.bidan_id, jk.tanggal_kontrol, jk.catatan, jk.status_notifikasi, jk.created_at, jk.updated_at,
		        p.nama_lengkap
		 FROM jadwal_kontrol jk
		 JOIN pasien p ON p.id = jk.pasien_id
		 ORDER BY jk.tanggal_kontrol ASC
		 LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.JadwalKontrol
	for rows.Next() {
		jk := models.JadwalKontrol{}
		var tgl time.Time
		err := rows.Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt, &jk.NamaPasien)
		if err != nil {
			return nil, 0, err
		}
		jk.TanggalKontrol = tgl
		list = append(list, jk)
	}
	return list, total, nil
}

func (r *JadwalRepository) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.JadwalKontrol, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jadwal_kontrol WHERE pasien_id = ?`, pasienID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, pasien_id, bidan_id, tanggal_kontrol, catatan, status_notifikasi, created_at, updated_at
		 FROM jadwal_kontrol WHERE pasien_id = ?
		 ORDER BY tanggal_kontrol ASC LIMIT ? OFFSET ?`, pasienID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.JadwalKontrol
	for rows.Next() {
		jk := models.JadwalKontrol{}
		var tgl time.Time
		err := rows.Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		jk.TanggalKontrol = tgl
		list = append(list, jk)
	}
	return list, total, nil
}

func (r *JadwalRepository) FindUpcomingForNotification(ctx context.Context, targetDate time.Time) ([]models.JadwalKontrol, error) {
	dateStr := targetDate.Format("2006-01-02")
	rows, err := r.db.QueryContext(ctx,
		`SELECT jk.id, jk.pasien_id, jk.bidan_id, jk.tanggal_kontrol, jk.catatan, jk.status_notifikasi, jk.created_at, jk.updated_at,
		        p.nama_lengkap, p.no_wa, p.jenis_kelamin
		 FROM jadwal_kontrol jk
		 JOIN pasien p ON p.id = jk.pasien_id
		 WHERE DATE(jk.tanggal_kontrol) = ? AND jk.status_notifikasi = 'belum'`, dateStr,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.JadwalKontrol
	for rows.Next() {
		jk := models.JadwalKontrol{}
		var tgl time.Time
		err := rows.Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt, &jk.NamaPasien, &jk.NoWaPasien, &jk.JenisKelaminPasien)
		if err != nil {
			return nil, err
		}
		jk.TanggalKontrol = tgl
		list = append(list, jk)
	}
	return list, nil
}

func (r *JadwalRepository) Update(ctx context.Context, id int, tanggal time.Time, catatan *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE jadwal_kontrol SET tanggal_kontrol=?, catatan=? WHERE id=?`, tanggal, catatan, id,
	)
	return err
}

func (r *JadwalRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM jadwal_kontrol WHERE id = ?`, id)
	return err
}

func (r *JadwalRepository) UpdateNotifStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE jadwal_kontrol SET status_notifikasi=? WHERE id=?`, status, id,
	)
	return err
}

func (r *JadwalRepository) FindTodaySchedules(ctx context.Context, dateStr string) ([]models.JadwalKontrol, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT jk.id, jk.pasien_id, jk.bidan_id, jk.tanggal_kontrol, jk.catatan, jk.status_notifikasi, jk.created_at, jk.updated_at,
		        p.nama_lengkap
		 FROM jadwal_kontrol jk
		 JOIN pasien p ON p.id = jk.pasien_id
		 WHERE DATE(jk.tanggal_kontrol) = ?
		 ORDER BY jk.tanggal_kontrol ASC`, dateStr,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.JadwalKontrol{}
	for rows.Next() {
		jk := models.JadwalKontrol{}
		var tgl time.Time
		err := rows.Scan(&jk.ID, &jk.PasienID, &jk.BidanID, &tgl, &jk.Catatan, &jk.StatusNotifikasi, &jk.CreatedAt, &jk.UpdatedAt, &jk.NamaPasien)
		if err != nil {
			return nil, err
		}
		jk.TanggalKontrol = tgl
		list = append(list, jk)
	}
	return list, nil
}

