package repositories

import (
	"context"
	"database/sql"
)

// KlinikRepository handles database operations for the klinik_status table.
type KlinikRepository struct {
	db *sql.DB
}

// NewKlinikRepository creates a new KlinikRepository instance.
func NewKlinikRepository(db *sql.DB) *KlinikRepository {
	return &KlinikRepository{db: db}
}

// GetStatus retrieves the latest clinic status and notes.
func (r *KlinikRepository) GetStatus(ctx context.Context) (status string, catatan string, err error) {
	var catatanNull sql.NullString
	err = r.db.QueryRowContext(ctx,
		`SELECT status, catatan FROM klinik_status ORDER BY updated_at DESC LIMIT 1`,
	).Scan(&status, &catatanNull)
	if err != nil {
		if err == sql.ErrNoRows {
			return "tutup", "", nil
		}
		return "tutup", "", err
	}
	if catatanNull.Valid {
		catatan = catatanNull.String
	}
	return status, catatan, nil
}

// SetStatus updates or inserts the clinic status for a given bidan.
func (r *KlinikRepository) SetStatus(ctx context.Context, bidanID int, status, catatan string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE klinik_status SET status = ?, catatan = ? WHERE bidan_id = ?`,
		status, catatan, bidanID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO klinik_status (bidan_id, status, catatan) VALUES (?, ?, ?)`,
			bidanID, status, catatan,
		)
		return err
	}
	return nil
}

// GetStatusString returns just the status string, defaulting to "tutup" on error.
func (r *KlinikRepository) GetStatusString(ctx context.Context) string {
	var status string
	err := r.db.QueryRowContext(ctx,
		`SELECT status FROM klinik_status ORDER BY updated_at DESC LIMIT 1`,
	).Scan(&status)
	if err != nil {
		return "tutup"
	}
	return status
}
