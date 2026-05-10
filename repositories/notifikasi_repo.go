package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

type NotifikasiRepository struct {
	db *sql.DB
}

func NewNotifikasiRepository(db *sql.DB) *NotifikasiRepository {
	return &NotifikasiRepository{db: db}
}

func (r *NotifikasiRepository) Create(ctx context.Context, pasienID int, jadwalKontrolID *int, judul, pesan, statusKirim string, sentAt *time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notifikasi (pasien_id, jadwal_kontrol_id, judul, pesan, status_kirim, sent_at)
		 VALUES (?,?,?,?,?,?)`,
		pasienID, jadwalKontrolID, judul, pesan, statusKirim, sentAt,
	)
	if err != nil {
		return fmt.Errorf("create notifikasi: %w", err)
	}
	return nil
}

func (r *NotifikasiRepository) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Notifikasi, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifikasi WHERE pasien_id = ?`, pasienID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, pasien_id, jadwal_kontrol_id, judul, pesan, channel, status_kirim, is_read, sent_at, created_at
		 FROM notifikasi WHERE pasien_id = ?
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`, pasienID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Notifikasi
	for rows.Next() {
		n := models.Notifikasi{}
		err := rows.Scan(&n.ID, &n.PasienID, &n.JadwalKontrolID, &n.Judul, &n.Pesan, &n.Channel, &n.StatusKirim, &n.IsRead, &n.SentAt, &n.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, n)
	}
	return list, total, nil
}

func (r *NotifikasiRepository) MarkAsRead(ctx context.Context, id, pasienID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifikasi SET is_read = 1 WHERE id = ? AND pasien_id = ?`, id, pasienID)
	return err
}
