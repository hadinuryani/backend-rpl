package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/dto"
	"ic-plus-backend/repositories"
)

type RekamMedisService struct {
	rmRepo      repositories.RekamMedisRepo
	antrianRepo repositories.AntrianRepo
	jadwalRepo  repositories.JadwalRepo
	pasienRepo  repositories.PasienRepo
	db          *sql.DB
}

func NewRekamMedisService(rmRepo repositories.RekamMedisRepo, antrianRepo repositories.AntrianRepo, jadwalRepo repositories.JadwalRepo, pasienRepo repositories.PasienRepo, db *sql.DB) *RekamMedisService {
	return &RekamMedisService{rmRepo: rmRepo, antrianRepo: antrianRepo, jadwalRepo: jadwalRepo, pasienRepo: pasienRepo, db: db}
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
		Jumlah      int
		Dosis       string
		AturanPakai string
		Catatan     *string
	}, len(req.Resep))
	for i, d := range req.Resep {
		details[i].ObatID = d.ObatID
		details[i].Jumlah = d.Jumlah
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

	// Retrieve antrian and pasien ID if we need to create schedule or update pregnancy status
	var pasienID int
	if (req.PerluKontrol && req.TanggalKontrol != "") || req.IsHamil != nil {
		antrian, err := s.antrianRepo.FindByID(ctx, req.AntrianID)
		if err != nil || antrian == nil {
			return 0, 0, fmt.Errorf("antrian tidak ditemukan")
		}
		pasienID = antrian.PasienID

		// Update pregnancy status if provided
		if req.IsHamil != nil {
			err = s.pasienRepo.UpdateIsHamilTx(ctx, tx, pasienID, *req.IsHamil)
			if err != nil {
				return 0, 0, fmt.Errorf("update status kehamilan: %w", err)
			}
		}
	}

	// Auto-create control schedule if requested
	if req.PerluKontrol && req.TanggalKontrol != "" {
		tanggalKontrol, err := time.Parse("2006-01-02", req.TanggalKontrol)
		if err != nil {
			return 0, 0, fmt.Errorf("format tanggal kontrol tidak valid")
		}

		var catatanKontrol *string
		defaultCatatan := "Kontrol Rutin"
		if req.CatatanKontrol != "" {
			catatanKontrol = &req.CatatanKontrol
		} else {
			catatanKontrol = &defaultCatatan
		}
		
		_, err = s.jadwalRepo.CreateTx(ctx, tx, pasienID, bidanID, tanggalKontrol, catatanKontrol)
		if err != nil {
			return 0, 0, fmt.Errorf("create jadwal kontrol: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("commit transaction: %w", err)
	}

	return rmID, resepID, nil
}
