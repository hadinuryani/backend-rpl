package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

// PasienRepository handles database operations for the pasien table.
type PasienRepository struct {
	db *sql.DB
}

// NewPasienRepository creates a new PasienRepository instance.
func NewPasienRepository(db *sql.DB) *PasienRepository {
	return &PasienRepository{db: db}
}

// Create inserts a new patient record.
func (r *PasienRepository) Create(ctx context.Context, userID int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) (*models.Pasien, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO pasien (user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, nama, tglLahir, jenisKelamin, alamat, noWa, golDarah,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pasien: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}
	return r.FindByID(ctx, int(id))
}

// CreateByBidan inserts a new patient record within a transaction.
func (r *PasienRepository) CreateByBidan(ctx context.Context, tx *sql.Tx, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string, userID int) (*models.Pasien, error) {
	result, err := tx.ExecContext(ctx,
		`INSERT INTO pasien (user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, nama, tglLahir, jenisKelamin, alamat, noWa, golDarah,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pasien by bidan: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	pasien := &models.Pasien{}
	err = tx.QueryRowContext(ctx,
		`SELECT id, user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah, created_at, updated_at,
		        TIMESTAMPDIFF(YEAR, tanggal_lahir, CURDATE()) as umur, is_hamil
		 FROM pasien WHERE id = ?`,
		id,
	).Scan(&pasien.ID, &pasien.UserID, &pasien.NamaLengkap, &pasien.TanggalLahir, &pasien.JenisKelamin,
		&pasien.Alamat, &pasien.NoWa, &pasien.GolonganDarah, &pasien.CreatedAt, &pasien.UpdatedAt, &pasien.Umur, &pasien.IsHamil)

	return pasien, err
}

// FindByUserID retrieves a patient by their user ID.
func (r *PasienRepository) FindByUserID(ctx context.Context, userID int) (*models.Pasien, error) {
	pasien := &models.Pasien{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah, created_at, updated_at,
		        TIMESTAMPDIFF(YEAR, tanggal_lahir, CURDATE()) as umur, is_hamil
		 FROM pasien WHERE user_id = ?`,
		userID,
	).Scan(&pasien.ID, &pasien.UserID, &pasien.NamaLengkap, &pasien.TanggalLahir, &pasien.JenisKelamin,
		&pasien.Alamat, &pasien.NoWa, &pasien.GolonganDarah, &pasien.CreatedAt, &pasien.UpdatedAt, &pasien.Umur, &pasien.IsHamil)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pasien by user ID: %w", err)
	}
	return pasien, nil
}

// FindByID retrieves a patient by their pasien ID.
func (r *PasienRepository) FindByID(ctx context.Context, id int) (*models.Pasien, error) {
	pasien := &models.Pasien{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah, created_at, updated_at,
		        TIMESTAMPDIFF(YEAR, tanggal_lahir, CURDATE()) as umur, is_hamil
		 FROM pasien WHERE id = ?`,
		id,
	).Scan(&pasien.ID, &pasien.UserID, &pasien.NamaLengkap, &pasien.TanggalLahir, &pasien.JenisKelamin,
		&pasien.Alamat, &pasien.NoWa, &pasien.GolonganDarah, &pasien.CreatedAt, &pasien.UpdatedAt, &pasien.Umur, &pasien.IsHamil)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pasien by ID: %w", err)
	}
	return pasien, nil
}

// Update updates a patient's profile data.
func (r *PasienRepository) Update(ctx context.Context, id int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pasien SET nama_lengkap=?, tanggal_lahir=?, jenis_kelamin=?, alamat=?, no_wa=?, golongan_darah=?
		 WHERE id = ?`,
		nama, tglLahir, jenisKelamin, alamat, noWa, golDarah, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update pasien: %w", err)
	}
	return nil
}

// FindAll retrieves all patients with optional search and pagination.
func (r *PasienRepository) FindAll(ctx context.Context, search string, limit, offset int) ([]models.Pasien, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM pasien WHERE (? = '' OR nama_lengkap LIKE CONCAT('%', ?, '%'))`
	err := r.db.QueryRowContext(ctx, countQuery, search, search).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pasien: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah, created_at, updated_at,
		        TIMESTAMPDIFF(YEAR, tanggal_lahir, CURDATE()) as umur, is_hamil
		 FROM pasien
		 WHERE (? = '' OR nama_lengkap LIKE CONCAT('%', ?, '%'))
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`,
		search, search, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pasien: %w", err)
	}
	defer rows.Close()

	var patients []models.Pasien
	for rows.Next() {
		p := models.Pasien{}
		err := rows.Scan(&p.ID, &p.UserID, &p.NamaLengkap, &p.TanggalLahir, &p.JenisKelamin,
			&p.Alamat, &p.NoWa, &p.GolonganDarah, &p.CreatedAt, &p.UpdatedAt, &p.Umur, &p.IsHamil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pasien row: %w", err)
		}
		patients = append(patients, p)
	}

	return patients, total, nil
}

func (r *PasienRepository) FindByNoWa(ctx context.Context, noWa string) (*models.Pasien, error) {
	pasien := &models.Pasien{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, nama_lengkap, tanggal_lahir, jenis_kelamin, alamat, no_wa, golongan_darah, created_at, updated_at,
		        TIMESTAMPDIFF(YEAR, tanggal_lahir, CURDATE()) as umur, is_hamil
		 FROM pasien WHERE no_wa = ?`,
		noWa,
	).Scan(&pasien.ID, &pasien.UserID, &pasien.NamaLengkap, &pasien.TanggalLahir, &pasien.JenisKelamin,
		&pasien.Alamat, &pasien.NoWa, &pasien.GolonganDarah, &pasien.CreatedAt, &pasien.UpdatedAt, &pasien.Umur, &pasien.IsHamil)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pasien by WhatsApp number: %w", err)
	}
	return pasien, nil
}

// UpdateIsHamil updates the patient's pregnancy status.
func (r *PasienRepository) UpdateIsHamil(ctx context.Context, id int, isHamil bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE pasien SET is_hamil = ? WHERE id = ?", isHamil, id)
	if err != nil {
		return fmt.Errorf("failed to update is_hamil: %w", err)
	}
	return nil
}

// UpdateIsHamilTx updates the patient's pregnancy status within a transaction.
func (r *PasienRepository) UpdateIsHamilTx(ctx context.Context, tx *sql.Tx, id int, isHamil bool) error {
	_, err := tx.ExecContext(ctx, "UPDATE pasien SET is_hamil = ? WHERE id = ?", isHamil, id)
	if err != nil {
		return fmt.Errorf("failed to update is_hamil in transaction: %w", err)
	}
	return nil
}

