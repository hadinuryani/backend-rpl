package services

import (
	"context"
	"testing"
)

func TestAntrianService_CreateAntrian(t *testing.T) {
	antrianRepo := &MockAntrianRepo{}
	klinikRepo := &MockKlinikRepo{Status: "buka"}
	service := NewAntrianService(antrianRepo, klinikRepo)

	ctx := context.Background()

	// Test Successful Queue Creation
	antrianRepo.CountTodayVal = 3
	antrian, err := service.CreateAntrian(ctx, 1, "2026-05-28", "Batuk dan pilek")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if antrian == nil {
		t.Fatal("expected antrian to be returned")
	}
	if antrian.NoAntrian == "" {
		t.Error("expected a queue number to be generated")
	}

	// Test Failed Queue Creation - Clinic Closed
	klinikRepo.Status = "tutup"
	_, err = service.CreateAntrian(ctx, 1, "2026-05-28", "Demam")
	if err == nil {
		t.Error("expected error when clinic is closed, got nil")
	}

	// Test Failed Queue Creation - Invalid Date Format
	klinikRepo.Status = "buka"
	_, err = service.CreateAntrian(ctx, 1, "invalid-date", "Sakit kepala")
	if err == nil {
		t.Error("expected error for invalid date format, got nil")
	}
}
