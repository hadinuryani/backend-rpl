package services

import (
	"context"
	"database/sql"
	"fmt"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
)

type RekamMedisService struct {
	rmRepo      *repositories.RekamMedisRepository
	antrianRepo *repositories.AntrianRepository
	db          *sql.DB
}

func NewRekamMedisService(rmRepo *repositories.RekamMedisRepository, antrianRepo *repositories.AntrianRepository, db *sql.DB) *RekamMedisService {
	return &RekamMedisService{rmRepo: rmRepo, antrianRepo: antrianRepo, db: db}
}

// CreateWithResep creates rekam medis + resep + detail_resep + marks antrian as selesai, all in one transaction.
func (s *RekamMedisService) CreateWithResep(ctx context.Context, bidanID int, req dto.CreateRekamMedisRequest) (int, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var tekananDarah *string
	if req.TekananDarah != "" {
		tekananDarah = &req.TekananDarah
	}
	var kondisiJanin *string
	if req.KondisiJanin != "" {
		kondisiJanin = &req.KondisiJanin
	}
	var catatanTambahan *string
	if req.CatatanTambahan != "" {
		catatanTambahan = &req.CatatanTambahan
	}

	details := make([]struct {
		ObatID      int
		Dosis       string
		AturanPakai string
		Catatan     *string
	}, len(req.Resep))
	for i, d := range req.Resep {
		details[i].ObatID = d.ObatID
		details[i].Dosis = d.Dosis
		details[i].AturanPakai = d.AturanPakai
		if d.Catatan != "" {
			cat := d.Catatan
			details[i].Catatan = &cat
		}
	}

	rmID, resepID, err := s.rmRepo.CreateWithResepTx(ctx, tx, bidanID, req.AntrianID, req.KeluhanUtama, tekananDarah, req.BeratBadan, req.TinggiFundusUteri, kondisiJanin, catatanTambahan, details)
	if err != nil {
		return 0, 0, err
	}

	err = s.antrianRepo.UpdateStatusTx(ctx, tx, req.AntrianID, "selesai")
	if err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("commit transaction: %w", err)
	}

	return rmID, resepID, nil
}
