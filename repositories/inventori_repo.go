package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

type InventoriRepository struct {
	db *sql.DB
}

func NewInventoriRepository(db *sql.DB) *InventoriRepository {
	return &InventoriRepository{db: db}
}

// --- Obat CRUD ---

func (r *InventoriRepository) CreateObat(ctx context.Context, nama, kategori, satuan string, stokMin int) (*models.Obat, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO obat (nama_obat, kategori, satuan, stok_minimum)
		 VALUES (?,?,?,?)`,
		nama, kategori, satuan, stokMin,
	)
	if err != nil {
		return nil, fmt.Errorf("create obat: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	o := &models.Obat{}
	err = r.db.QueryRowContext(ctx,
		`SELECT id, nama_obat, kategori, satuan, stok_minimum, created_at, updated_at FROM obat WHERE id=?`, id,
	).Scan(&o.ID, &o.NamaObat, &o.Kategori, &o.Satuan, &o.StokMinimum, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *InventoriRepository) CreateObatWithInventori(ctx context.Context, nama, kategori, satuan string, stokMin, jumlahStok int, tanggalKadaluarsa *time.Time, batchNumber *string) (*models.Obat, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`INSERT INTO obat (nama_obat, kategori, satuan, stok_minimum)
		 VALUES (?,?,?,?)`,
		nama, kategori, satuan, stokMin,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	statusStok := calculateStatusStok(jumlahStok, stokMin, tanggalKadaluarsa)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO inventori (obat_id, jumlah_stok, tanggal_kadaluarsa, batch_number, status_stok)
		 VALUES (?,?,?,?,?)`,
		id, jumlahStok, tanggalKadaluarsa, batchNumber, statusStok,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	o := &models.Obat{}
	err = r.db.QueryRowContext(ctx,
		`SELECT id, nama_obat, kategori, satuan, stok_minimum, created_at, updated_at FROM obat WHERE id=?`, id,
	).Scan(&o.ID, &o.NamaObat, &o.Kategori, &o.Satuan, &o.StokMinimum, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *InventoriRepository) FindAllObat(ctx context.Context, limit, offset int) ([]models.Obat, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM obat`).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, nama_obat, kategori, satuan, stok_minimum, created_at, updated_at FROM obat ORDER BY nama_obat ASC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Obat
	for rows.Next() {
		o := models.Obat{}
		rows.Scan(&o.ID, &o.NamaObat, &o.Kategori, &o.Satuan, &o.StokMinimum, &o.CreatedAt, &o.UpdatedAt)
		list = append(list, o)
	}
	return list, total, nil
}

func (r *InventoriRepository) UpdateObat(ctx context.Context, id int, nama, kategori, satuan string, stokMin *int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE obat SET nama_obat=?, kategori=?, satuan=?, stok_minimum=COALESCE(?, stok_minimum) WHERE id=?`,
		nama, kategori, satuan, stokMin, id)
	return err
}

func (r *InventoriRepository) DeleteObat(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM obat WHERE id = ?`, id)
	return err
}

// --- Inventori ---

func (r *InventoriRepository) FindAllInventori(ctx context.Context, statusFilter string, limit, offset int) ([]models.Inventori, int, error) {
	var total int
	countQ := `SELECT COUNT(*) FROM inventori i JOIN obat o ON o.id = i.obat_id`
	countArgs := []interface{}{}
	if statusFilter != "" {
		countQ += ` WHERE i.status_stok = ?`
		countArgs = append(countArgs, statusFilter)
	}
	r.db.QueryRowContext(ctx, countQ, countArgs...).Scan(&total)

	query := `SELECT i.id, i.obat_id, i.jumlah_stok, i.tanggal_kadaluarsa, i.batch_number, i.status_stok, i.updated_at,
	                  o.nama_obat, o.kategori, o.satuan, o.stok_minimum
	           FROM inventori i JOIN obat o ON o.id = i.obat_id`
	args := []interface{}{}
	if statusFilter != "" {
		query += ` WHERE i.status_stok = ? ORDER BY o.nama_obat ASC LIMIT ? OFFSET ?`
		args = append(args, statusFilter, limit, offset)
	} else {
		query += ` ORDER BY o.nama_obat ASC LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Inventori
	for rows.Next() {
		inv := models.Inventori{}
		var tgl []byte
		rows.Scan(&inv.ID, &inv.ObatID, &inv.JumlahStok, &tgl, &inv.BatchNumber, &inv.StatusStok, &inv.UpdatedAt,
			&inv.NamaObat, &inv.Kategori, &inv.Satuan, &inv.StokMinimum)
		if len(tgl) > 0 {
			parsed, _ := time.Parse("2006-01-02", string(tgl))
			inv.TanggalKadaluarsa = &parsed
		}
		list = append(list, inv)
	}
	return list, total, nil
}

func (r *InventoriRepository) StokMasuk(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan *string, tanggalKadaluarsa *time.Time, batchNumber *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update inventory
	updateQ := `UPDATE inventori SET jumlah_stok = jumlah_stok + ?`
	args := []interface{}{jumlah}
	
	if tanggalKadaluarsa != nil {
		updateQ += `, tanggal_kadaluarsa = ?`
		args = append(args, tanggalKadaluarsa)
	}
	if batchNumber != nil {
		updateQ += `, batch_number = ?`
		args = append(args, batchNumber)
	}
	updateQ += ` WHERE id = ?`
	args = append(args, inventoriID)

	_, err = tx.ExecContext(ctx, updateQ, args...)
	if err != nil {
		return err
	}

	// Insert history
	_, err = tx.ExecContext(ctx,
		`INSERT INTO riwayat_stok (inventori_id, bidan_id, jenis_transaksi, jumlah, keterangan) VALUES (?,?,'masuk',?,?)`,
		inventoriID, bidanID, jumlah, keterangan)
	if err != nil {
		return err
	}

	// Recalculate status
	if err := r.recalcStatusTx(ctx, tx, inventoriID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *InventoriRepository) StokKeluar(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check current stock
	var currentStok int
	err = tx.QueryRowContext(ctx, `SELECT jumlah_stok FROM inventori WHERE id = ? FOR UPDATE`, inventoriID).Scan(&currentStok)
	if err != nil {
		return fmt.Errorf("inventory not found: %w", err)
	}
	if currentStok < jumlah {
		return fmt.Errorf("stok tidak mencukupi (tersedia: %d, diminta: %d)", currentStok, jumlah)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE inventori SET jumlah_stok = jumlah_stok - ? WHERE id = ?`,
		jumlah, inventoriID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO riwayat_stok (inventori_id, bidan_id, jenis_transaksi, jumlah, keterangan) VALUES (?,?,'keluar',?,?)`,
		inventoriID, bidanID, jumlah, keterangan)
	if err != nil {
		return err
	}

	if err := r.recalcStatusTx(ctx, tx, inventoriID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *InventoriRepository) GetRiwayat(ctx context.Context, limit, offset int) ([]models.RiwayatStok, int, error) {
	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM riwayat_stok`).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		`SELECT rs.id, rs.inventori_id, rs.bidan_id, rs.jenis_transaksi, rs.jumlah, rs.keterangan, rs.created_at,
		        o.nama_obat, b.nama_lengkap
		 FROM riwayat_stok rs
		 JOIN inventori i ON i.id = rs.inventori_id
		 JOIN obat o ON o.id = i.obat_id
		 JOIN bidan b ON b.id = rs.bidan_id
		 ORDER BY rs.created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.RiwayatStok
	for rows.Next() {
		rs := models.RiwayatStok{}
		rows.Scan(&rs.ID, &rs.InventoriID, &rs.BidanID, &rs.JenisTransaksi, &rs.Jumlah, &rs.Keterangan, &rs.CreatedAt, &rs.NamaObat, &rs.NamaBidan)
		list = append(list, rs)
	}
	return list, total, nil
}

func (r *InventoriRepository) CountCritical(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM inventori WHERE status_stok IN ('hampir_habis','habis','kadaluarsa','hampir_kadaluarsa')`).Scan(&count)
	return count, err
}

func (r *InventoriRepository) recalcStatusTx(ctx context.Context, tx *sql.Tx, inventoriID int) error {
	var jumlahStok, stokMin int
	var tgl []byte
	var tanggalKadaluarsa *time.Time
	
	err := tx.QueryRowContext(ctx,
		`SELECT i.jumlah_stok, o.stok_minimum, i.tanggal_kadaluarsa
		 FROM inventori i JOIN obat o ON o.id = i.obat_id WHERE i.id = ?`, inventoriID,
	).Scan(&jumlahStok, &stokMin, &tgl)
	if err != nil {
		return err
	}
	if len(tgl) > 0 {
		parsed, _ := time.Parse("2006-01-02", string(tgl))
		tanggalKadaluarsa = &parsed
	}

	status := calculateStatusStok(jumlahStok, stokMin, tanggalKadaluarsa)
	_, err = tx.ExecContext(ctx, `UPDATE inventori SET status_stok = ? WHERE id = ?`, status, inventoriID)
	return err
}

func calculateStatusStok(jumlahStok, stokMinimum int, tanggalKadaluarsa *time.Time) string {
	if jumlahStok == 0 {
		return "habis"
	}
	if tanggalKadaluarsa != nil {
		now := time.Now()
		if !tanggalKadaluarsa.After(now) {
			return "kadaluarsa"
		}
		if !tanggalKadaluarsa.After(now.AddDate(0, 0, 30)) {
			return "hampir_kadaluarsa"
		}
	}
	if jumlahStok <= stokMinimum {
		return "hampir_habis"
	}
	return "aman"
}
