package services

import (
	"context"
	"time"

	"ic-plus-backend/repositories"
)

type NotifikasiService struct {
	repo repositories.NotifikasiRepo
}

func NewNotifikasiService(repo repositories.NotifikasiRepo) *NotifikasiService {
	return &NotifikasiService{repo: repo}
}

func (s *NotifikasiService) Create(ctx context.Context, pasienID int, jadwalKontrolID *int, judul, pesan, statusKirim string) error {
	now := time.Now()
	return s.repo.Create(ctx, pasienID, jadwalKontrolID, judul, pesan, statusKirim, &now)
}

func (s *NotifikasiService) MarkAsRead(ctx context.Context, id, pasienID int) error {
	return s.repo.MarkAsRead(ctx, id, pasienID)
}
