package services

import (
	"context"
	"fmt"
	"time"

	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
)

type JadwalService struct {
	repo *repositories.JadwalRepository
}

func NewJadwalService(repo *repositories.JadwalRepository) *JadwalService {
	return &JadwalService{repo: repo}
}

func (s *JadwalService) Create(ctx context.Context, pasienID, bidanID int, tanggalStr, catatan string) (*models.JadwalKontrol, error) {
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}
	var catatanPtr *string
	if catatan != "" {
		catatanPtr = &catatan
	}
	return s.repo.Create(ctx, pasienID, bidanID, tanggal, catatanPtr)
}

func (s *JadwalService) Update(ctx context.Context, id int, tanggalStr, catatan string) error {
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return fmt.Errorf("format tanggal tidak valid")
	}
	var catatanPtr *string
	if catatan != "" {
		catatanPtr = &catatan
	}
	return s.repo.Update(ctx, id, tanggal, catatanPtr)
}

func (s *JadwalService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
