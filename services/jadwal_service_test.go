package services

import (
	"context"
	"testing"
)

func TestJadwalService_CRUD(t *testing.T) {
	repo := NewMockJadwalRepo()
	service := NewJadwalService(repo)

	ctx := context.Background()

	// Test Successful Creation
	jk, err := service.Create(ctx, 1, 2, "2026-06-15", "Kontrol Rutin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if jk.PasienID != 1 || jk.BidanID != 2 || *jk.Catatan != "Kontrol Rutin" {
		t.Errorf("created schedule properties mismatched: %+v", jk)
	}

	// Test Failed Creation - Invalid Date Format
	_, err = service.Create(ctx, 1, 2, "invalid-date-format", "")
	if err == nil {
		t.Error("expected error for invalid date format in Create, got nil")
	}

	// Test Successful Update
	err = service.Update(ctx, jk.ID, "2026-06-20", "Tanggal diubah")
	if err != nil {
		t.Fatalf("expected no error during update, got %v", err)
	}
	updatedJk := repo.Schedules[jk.ID]
	if updatedJk.TanggalKontrol.Format("2006-01-02") != "2026-06-20" || *updatedJk.Catatan != "Tanggal diubah" {
		t.Errorf("schedule was not updated correctly: %+v", updatedJk)
	}

	// Test Failed Update - Invalid Date Format
	err = service.Update(ctx, jk.ID, "invalid-date", "")
	if err == nil {
		t.Error("expected error for invalid date format in Update, got nil")
	}

	// Test Successful Delete
	err = service.Delete(ctx, jk.ID)
	if err != nil {
		t.Fatalf("expected no error during delete, got %v", err)
	}
	if _, exists := repo.Schedules[jk.ID]; exists {
		t.Error("expected schedule to be deleted from repo")
	}
}
