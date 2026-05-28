package services

import (
	"context"
	"fmt"
	"time"

	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"
)

type AntrianService struct {
	antrianRepo repositories.AntrianRepo
	klinikRepo  repositories.KlinikRepo
}

func NewAntrianService(antrianRepo repositories.AntrianRepo, klinikRepo repositories.KlinikRepo) *AntrianService {
	return &AntrianService{antrianRepo: antrianRepo, klinikRepo: klinikRepo}
}

// CreateAntrian registers a new visit. Auto-generates queue number, checks clinic status.
func (s *AntrianService) CreateAntrian(ctx context.Context, pasienID int, tanggalStr, keluhan string) (*models.Antrian, error) {
	klinikStatus := s.klinikRepo.GetStatusString(ctx)
	if klinikStatus != "buka" {
		return nil, fmt.Errorf("klinik sedang tutup, pendaftaran tidak tersedia")
	}

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}

	count, err := s.antrianRepo.CountToday(ctx, tanggalStr)
	if err != nil {
		return nil, err
	}

	noAntrian := utils.GenerateNoAntrian(count, tanggal.Weekday())
	return s.antrianRepo.Create(ctx, pasienID, tanggal, noAntrian, keluhan)
}
