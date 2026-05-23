package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

type RekamMedisRepository struct {
	db *sql.DB
}

func NewRekamMedisRepository(db *sql.DB) *RekamMedisRepository {
	return &RekamMedisRepository{db: db}
}

// CreateWithResepTx creates rekam medis + resep + detail_resep in a transaction.
// Also deducts medicine stock from inventori for each prescribed item.
// Returns the created rekam medis ID and resep ID.
func (r *RekamMedisRepository) CreateWithResepTx(ctx context.Context, tx *sql.Tx, bidanID, antrianID int, keluhanUtama string, tekananDarah *string, beratBadan, tinggiFundus *float64, kondisiJanin, catatan *string, details []struct{ ObatID, Jumlah int; Dosis, AturanPakai string; Catatan *string }) (int, int, error) {
	// 1. Insert rekam_medis
	result, err := tx.ExecContext(ctx,
		`INSERT INTO rekam_medis (antrian_id, bidan_id, keluhan_utama, tekanan_darah, berat_badan, tinggi_fundus_uteri, kondisi_janin, catatan_tambahan)
		 VALUES (?,?,?,?,?,?,?,?)`,
		antrianID, bidanID, keluhanUtama, tekananDarah, beratBadan, tinggiFundus, kondisiJanin, catatan,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert rekam_medis: %w", err)
	}
	rmID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	// 2. Insert resep
	result, err = tx.ExecContext(ctx,
		`INSERT INTO resep (rekam_medis_id) VALUES (?)`, rmID,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert resep: %w", err)
	}
	resepID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	// 3. Insert detail_resep + deduct stock
	for _, d := range details {
		jumlah := d.Jumlah
		if jumlah < 1 {
			jumlah = 1
		}

		_, err = tx.ExecContext(ctx,
			`INSERT INTO detail_resep (resep_id, obat_id, jumlah, dosis, aturan_pakai, catatan)
			 VALUES (?,?,?,?,?,?)`,
			resepID, d.ObatID, jumlah, d.Dosis, d.AturanPakai, d.Catatan,
		)
		if err != nil {
			return 0, 0, fmt.Errorf("insert detail_resep: %w", err)
		}

		// 4. Deduct stock from inventori (by obat_id)
		var inventoriID, currentStok int
		err = tx.QueryRowContext(ctx,
			`SELECT id, jumlah_stok FROM inventori WHERE obat_id = ? FOR UPDATE`, d.ObatID,
		).Scan(&inventoriID, &currentStok)
		if err != nil {
			// No inventory record for this obat — skip deduction
			continue
		}
		if currentStok < jumlah {
			return 0, 0, fmt.Errorf("stok obat tidak mencukupi (tersedia: %d, dibutuhkan: %d)", currentStok, jumlah)
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE inventori SET jumlah_stok = jumlah_stok - ? WHERE id = ?`,
			jumlah, inventoriID)
		if err != nil {
			return 0, 0, fmt.Errorf("deduct stock: %w", err)
		}

		// 5. Log stock history
		keterangan := fmt.Sprintf("Resep pasien (antrian #%d)", antrianID)
		_, err = tx.ExecContext(ctx,
			`INSERT INTO riwayat_stok (inventori_id, bidan_id, jenis_transaksi, jumlah, keterangan)
			 VALUES (?,?,'keluar',?,?)`,
			inventoriID, bidanID, jumlah, keterangan)
		if err != nil {
			return 0, 0, fmt.Errorf("insert riwayat_stok: %w", err)
		}

		// 6. Recalculate stock status
		r.recalcStatusInTx(ctx, tx, inventoriID)
	}

	return int(rmID), int(resepID), nil
}

func (r *RekamMedisRepository) FindByID(ctx context.Context, id int) (*models.RekamMedis, error) {
	rm := &models.RekamMedis{}
	var tgl []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT rm.id, rm.antrian_id, rm.bidan_id, rm.keluhan_utama, rm.tekanan_darah, rm.berat_badan, rm.tinggi_fundus_uteri, rm.kondisi_janin, rm.catatan_tambahan, rm.created_at,
		        p.nama_lengkap, a.tanggal_kunjungan, b.nama_lengkap
		 FROM rekam_medis rm
		 JOIN antrian a ON a.id = rm.antrian_id
		 JOIN pasien p ON p.id = a.pasien_id
		 JOIN bidan b ON b.id = rm.bidan_id
		 WHERE rm.id = ?`, id,
	).Scan(&rm.ID, &rm.AntrianID, &rm.BidanID, &rm.KeluhanUtama, &rm.TekananDarah, &rm.BeratBadan, &rm.TinggiFundusUteri, &rm.KondisiJanin, &rm.CatatanTambahan, &rm.CreatedAt, &rm.NamaPasien, &tgl, &rm.NamaBidan)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find rekam_medis: %w", err)
	}
	if len(tgl) > 0 {
		parsed, _ := time.Parse("2006-01-02", string(tgl))
		rm.TanggalKunjungan = parsed.Format("2006-01-02")
	}
	return rm, nil
}

func (r *RekamMedisRepository) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.RekamMedis, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM rekam_medis rm JOIN antrian a ON a.id = rm.antrian_id WHERE a.pasien_id = ?`, pasienID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT rm.id, rm.antrian_id, rm.bidan_id, rm.keluhan_utama, rm.tekanan_darah, rm.berat_badan, rm.tinggi_fundus_uteri, rm.kondisi_janin, rm.catatan_tambahan, rm.created_at,
		        p.nama_lengkap, CAST(a.tanggal_kunjungan AS CHAR), b.nama_lengkap
		 FROM rekam_medis rm
		 JOIN antrian a ON a.id = rm.antrian_id
		 JOIN pasien p ON p.id = a.pasien_id
		 JOIN bidan b ON b.id = rm.bidan_id
		 WHERE a.pasien_id = ?
		 ORDER BY rm.created_at DESC LIMIT ? OFFSET ?`, pasienID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.RekamMedis
	for rows.Next() {
		rm := models.RekamMedis{}
		err := rows.Scan(&rm.ID, &rm.AntrianID, &rm.BidanID, &rm.KeluhanUtama, &rm.TekananDarah, &rm.BeratBadan, &rm.TinggiFundusUteri, &rm.KondisiJanin, &rm.CatatanTambahan, &rm.CreatedAt, &rm.NamaPasien, &rm.TanggalKunjungan, &rm.NamaBidan)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, rm)
	}
	return list, total, nil
}

func (r *RekamMedisRepository) Update(ctx context.Context, id int, keluhanUtama string, tekananDarah *string, beratBadan, tinggiFundus *float64, kondisiJanin, catatan *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE rekam_medis SET keluhan_utama=?, tekanan_darah=?, berat_badan=?, tinggi_fundus_uteri=?, kondisi_janin=?, catatan_tambahan=? WHERE id=?`,
		keluhanUtama, tekananDarah, beratBadan, tinggiFundus, kondisiJanin, catatan, id,
	)
	return err
}

func (r *RekamMedisRepository) FindResepByRekamMedisID(ctx context.Context, rmID int) (*models.Resep, error) {
	resep := &models.Resep{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, rekam_medis_id, created_at FROM resep WHERE rekam_medis_id = ?`, rmID,
	).Scan(&resep.ID, &resep.RekamMedisID, &resep.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT dr.id, dr.resep_id, dr.obat_id, dr.jumlah, dr.dosis, dr.aturan_pakai, dr.catatan, o.nama_obat
		 FROM detail_resep dr JOIN obat o ON o.id = dr.obat_id
		 WHERE dr.resep_id = ?`, resep.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := models.DetailResep{}
		err := rows.Scan(&d.ID, &d.ResepID, &d.ObatID, &d.Jumlah, &d.Dosis, &d.AturanPakai, &d.Catatan, &d.NamaObat)
		if err != nil {
			return nil, err
		}
		resep.Details = append(resep.Details, d)
	}
	return resep, nil
}

func (r *RekamMedisRepository) recalcStatusInTx(ctx context.Context, tx *sql.Tx, inventoriID int) error {
	var jumlahStok, stokMin int
	var tanggalKadaluarsa *time.Time
	
	err := tx.QueryRowContext(ctx,
		`SELECT i.jumlah_stok, o.stok_minimum, i.tanggal_kadaluarsa
		 FROM inventori i JOIN obat o ON o.id = i.obat_id WHERE i.id = ?`, inventoriID,
	).Scan(&jumlahStok, &stokMin, &tanggalKadaluarsa)
	if err != nil {
		return err
	}

	status := calculateStatusStok(jumlahStok, stokMin, tanggalKadaluarsa)
	_, err = tx.ExecContext(ctx, `UPDATE inventori SET status_stok = ? WHERE id = ?`, status, inventoriID)
	return err
}
