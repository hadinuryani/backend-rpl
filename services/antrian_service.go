package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"
)

type AntrianService struct {
	antrianRepo *repositories.AntrianRepository
	db          *sql.DB
}

func NewAntrianService(antrianRepo *repositories.AntrianRepository, db *sql.DB) *AntrianService {
	return &AntrianService{antrianRepo: antrianRepo, db: db}
}

// CreateAntrian registers a new visit. Auto-generates queue number, checks clinic status.
func (s *AntrianService) CreateAntrian(ctx context.Context, pasienID int, tanggalStr, keluhan string) (*models.Antrian, error) {
	var klinikStatus string
	err := s.db.QueryRowContext(ctx,
		`SELECT status FROM klinik_status ORDER BY updated_at DESC LIMIT 1`,
	).Scan(&klinikStatus)
	if err != nil {
		klinikStatus = "tutup"
	}
	if klinikStatus != "buka" {
		return nil, fmt.Errorf("klinik sedang tutup, pendaftaran tidak tersedia")
	}

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}

	count, err := s.antrianRepo.CountToday(ctx, tanggal)
	if err != nil {
		return nil, err
	}

	noAntrian := utils.GenerateNoAntrian(count)
	return s.antrianRepo.Create(ctx, pasienID, tanggal, noAntrian, keluhan)
}
