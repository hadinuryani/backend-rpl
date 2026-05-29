package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"ic-plus-backend/models"
)

// MockWAGateway implements WAGateway for testing
type MockWAGateway struct {
	SentMessages []struct {
		Target  string
		Message string
	}
	ErrorToReturn error
}

func (m *MockWAGateway) SendMessage(target, message string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	m.SentMessages = append(m.SentMessages, struct {
		Target  string
		Message string
	}{Target: target, Message: message})
	return nil
}

func TestSchedulerService_RunH1NotificationJob_Success(t *testing.T) {
	jadwalRepo := NewMockJadwalRepo()
	notifRepo := NewMockNotifikasiRepo()
	pengaturanRepo := NewMockPengaturanRepo()
	waGateway := &MockWAGateway{}

	// Setup settings
	pengaturanRepo.Settings["nama_klinik"] = "Klinik Indah Care Plus (IC+)"
	pengaturanRepo.Settings["jam_kontrol"] = "08:00 - selesai"

	// Setup tomorrow's schedule
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)

	// Patient details
	catatan := "Kontrol Rutin"
	jk := &models.JadwalKontrol{
		ID:               1,
		PasienID:         10,
		BidanID:          2,
		TanggalKontrol:   tomorrow,
		Catatan:          &catatan,
		StatusNotifikasi: "belum",
		NamaPasien:       "Doeng",
		NoWaPasien:       "081259277769",
	}
	jadwalRepo.Schedules[jk.ID] = jk

	service := NewSchedulerService(jadwalRepo, notifRepo, pengaturanRepo, waGateway)

	// Run job
	service.runH1NotificationJob()

	// Verify WA Gateway was called
	if len(waGateway.SentMessages) != 1 {
		t.Fatalf("expected 1 WA message to be sent, got %d", len(waGateway.SentMessages))
	}
	msg := waGateway.SentMessages[0]
	if msg.Target != "081259277769" {
		t.Errorf("expected target 081259277769, got %s", msg.Target)
	}
	if !strings.Contains(msg.Message, "Doeng") {
		t.Errorf("expected message to contain patient name Doeng, got: %s", msg.Message)
	}
	if !strings.Contains(msg.Message, "Klinik Indah Care Plus (IC+)") {
		t.Errorf("expected message to contain clinic name, got: %s", msg.Message)
	}

	// Verify notification record in DB
	if len(notifRepo.Notifs) != 1 {
		t.Fatalf("expected 1 notification record created, got %d", len(notifRepo.Notifs))
	}
	notif := notifRepo.Notifs[0]
	if notif.PasienID != 10 || *notif.JadwalKontrolID != 1 || notif.StatusKirim != "terkirim" {
		t.Errorf("unexpected notification properties: %+v", notif)
	}

	// Verify status in schedule repo updated to "terkirim"
	updatedJk := jadwalRepo.Schedules[1]
	if updatedJk.StatusNotifikasi != "terkirim" {
		t.Errorf("expected schedule status to be 'terkirim', got '%s'", updatedJk.StatusNotifikasi)
	}
}

func TestSchedulerService_RunH1NotificationJob_GatewayError(t *testing.T) {
	jadwalRepo := NewMockJadwalRepo()
	notifRepo := NewMockNotifikasiRepo()
	pengaturanRepo := NewMockPengaturanRepo()
	waGateway := &MockWAGateway{
		ErrorToReturn: errors.New("gateway offline"),
	}

	// Setup tomorrow's schedule
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)

	catatan := "Kontrol Rutin"
	jk := &models.JadwalKontrol{
		ID:               1,
		PasienID:         10,
		BidanID:          2,
		TanggalKontrol:   tomorrow,
		Catatan:          &catatan,
		StatusNotifikasi: "belum",
		NamaPasien:       "Doeng",
		NoWaPasien:       "081259277769",
	}
	jadwalRepo.Schedules[jk.ID] = jk

	service := NewSchedulerService(jadwalRepo, notifRepo, pengaturanRepo, waGateway)

	// Run job
	service.runH1NotificationJob()

	// Verify notification record in DB is saved as "gagal"
	if len(notifRepo.Notifs) != 1 {
		t.Fatalf("expected 1 notification record created, got %d", len(notifRepo.Notifs))
	}
	notif := notifRepo.Notifs[0]
	if notif.StatusKirim != "gagal" {
		t.Errorf("expected notification status to be 'gagal', got '%s'", notif.StatusKirim)
	}

	// Verify schedule notification status updated to "gagal"
	updatedJk := jadwalRepo.Schedules[1]
	if updatedJk.StatusNotifikasi != "gagal" {
		t.Errorf("expected schedule status to be 'gagal', got '%s'", updatedJk.StatusNotifikasi)
	}
}
