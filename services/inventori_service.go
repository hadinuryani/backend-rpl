package services

import (
	"context"
	"fmt"
	"time"

	"ic-plus-backend/repositories"
)

type InventoriService struct {
	repo repositories.InventoriRepo
}

func NewInventoriService(repo repositories.InventoriRepo) *InventoriService {
	return &InventoriService{repo: repo}
}

func (s *InventoriService) StokMasuk(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan, tanggalKadaluarsaStr, batchNumber string) error {
	var keteranganPtr *string
	if keterangan != "" {
		keteranganPtr = &keterangan
	}
	var tanggalKadaluarsa *time.Time
	if tanggalKadaluarsaStr != "" {
		t, err := time.Parse("2006-01-02", tanggalKadaluarsaStr)
		if err != nil {
			return fmt.Errorf("format tanggal kadaluarsa tidak valid")
		}
		tanggalKadaluarsa = &t
	}
	var batchPtr *string
	if batchNumber != "" {
		batchPtr = &batchNumber
	}
	return s.repo.StokMasuk(ctx, inventoriID, bidanID, jumlah, keteranganPtr, tanggalKadaluarsa, batchPtr)
}

func (s *InventoriService) StokKeluar(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan string) error {
	var keteranganPtr *string
	if keterangan != "" {
		keteranganPtr = &keterangan
	}
	return s.repo.StokKeluar(ctx, inventoriID, bidanID, jumlah, keteranganPtr)
}
