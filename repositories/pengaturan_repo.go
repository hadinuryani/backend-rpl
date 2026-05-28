package repositories

import (
	"context"
	"database/sql"
	"fmt"
)

// PengaturanRepository handles database operations for the pengaturan (settings) table.
type PengaturanRepository struct {
	db *sql.DB
}

// NewPengaturanRepository creates a new PengaturanRepository instance.
func NewPengaturanRepository(db *sql.DB) *PengaturanRepository {
	return &PengaturanRepository{db: db}
}

// GetAll retrieves all settings as a key-value map.
func (r *PengaturanRepository) GetAll(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT kunci, nilai FROM pengaturan`)
	if err != nil {
		return nil, fmt.Errorf("query pengaturan: %w", err)
	}
	defer rows.Close()

	settings := map[string]string{
		"waktu_pengingat": "08:00",
		"nama_klinik":     "Klinik Indah Care Plus (IC+)",
		"alamat_klinik":   "Jl. Indah Care No. 45, Jakarta",
		"jam_kontrol":     "08:00 - selesai",
	}

	for rows.Next() {
		var kunci, nilai string
		if err := rows.Scan(&kunci, &nilai); err == nil {
			settings[kunci] = nilai
		}
	}
	return settings, nil
}

// Upsert inserts or updates a setting by key.
func (r *PengaturanRepository) Upsert(ctx context.Context, kunci, nilai string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pengaturan (kunci, nilai) VALUES (?, ?) ON DUPLICATE KEY UPDATE nilai = ?`,
		kunci, nilai, nilai,
	)
	if err != nil {
		return fmt.Errorf("upsert pengaturan %s: %w", kunci, err)
	}
	return nil
}

// GetByKey retrieves a single setting value by key, returning fallback if not found.
func (r *PengaturanRepository) GetByKey(ctx context.Context, kunci string) (string, error) {
	var nilai string
	err := r.db.QueryRowContext(ctx, `SELECT nilai FROM pengaturan WHERE kunci = ?`, kunci).Scan(&nilai)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("query pengaturan %s: %w", kunci, err)
	}
	return nilai, nil
}
