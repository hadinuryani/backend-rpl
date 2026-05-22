package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"ic-plus-backend/repositories"

	"github.com/robfig/cron/v3"
)

// WAGateway is the interface for sending WhatsApp messages.
// Implement this interface to swap in a real sender.
type WAGateway interface {
	SendMessage(target, message string) error
}

// StubWAGateway logs messages to console instead of sending real WhatsApp messages.
type StubWAGateway struct{}

func (s *StubWAGateway) SendMessage(target, message string) error {
	log.Printf("📱 [WA STUB] To: %s\n   Message: %s\n", target, message)
	return nil
}

// SchedulerService manages the H-1 notification cron job.
type SchedulerService struct {
	jadwalRepo *repositories.JadwalRepository
	notifRepo  *repositories.NotifikasiRepository
	waGateway  WAGateway
	cron       *cron.Cron
	db         *sql.DB
	mu         sync.Mutex
	entryID    cron.EntryID
}

func NewSchedulerService(jadwalRepo *repositories.JadwalRepository, notifRepo *repositories.NotifikasiRepository, waGateway WAGateway, db *sql.DB) *SchedulerService {
	// Use WIB timezone (UTC+7)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	c := cron.New(cron.WithLocation(loc))
	return &SchedulerService{
		jadwalRepo: jadwalRepo,
		notifRepo:  notifRepo,
		waGateway:  waGateway,
		cron:       c,
		db:         db,
	}
}

// Start begins the cron scheduler with the H-1 notification job.
func (s *SchedulerService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read waktu_pengingat from database
	var waktu string
	err := s.db.QueryRow(`SELECT nilai FROM pengaturan WHERE kunci = 'waktu_pengingat'`).Scan(&waktu)
	if err != nil || waktu == "" {
		waktu = "08:00"
	}

	var hour, minute int
	_, _ = fmt.Sscanf(waktu, "%d:%d", &hour, &minute)

	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)
	s.cron.Start()

	id, err := s.cron.AddFunc(cronExpr, func() {
		log.Println("Running H-1 notification job...")
		s.runH1NotificationJob()
	})
	if err == nil {
		s.entryID = id
		log.Printf(" Scheduler started (H-1 notifications daily at %02d:%02d WIB)", hour, minute)
	} else {
		log.Printf(" Failed to start scheduler: %v", err)
	}
}

// Reschedule dynamic updates the cron schedule of the H-1 notification job.
func (s *SchedulerService) Reschedule(timeStr string) error {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.entryID != 0 {
		s.cron.Remove(s.entryID)
	}

	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)
	id, err := s.cron.AddFunc(cronExpr, func() {
		log.Println("Running H-1 notification job...")
		s.runH1NotificationJob()
	})
	if err != nil {
		return fmt.Errorf("failed to schedule cron: %w", err)
	}
	s.entryID = id
	log.Printf(" Scheduler rescheduled successfully for daily at %02d:%02d WIB (expr: %s)", hour, minute, cronExpr)
	return nil
}

// Stop gracefully stops the cron scheduler.
func (s *SchedulerService) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println(" Scheduler stopped")
}

func (s *SchedulerService) runH1NotificationJob() {
	ctx := context.Background()
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)

	schedules, err := s.jadwalRepo.FindUpcomingForNotification(ctx, tomorrow)
	if err != nil {
		log.Printf("H-1 job error: %v\n", err)
		return
	}

	if len(schedules) == 0 {
		log.Println(" No upcoming schedules for H-1 notification")
		return
	}

	// Load clinic settings
	namaKlinik := "Klinik Indah Care Plus (IC+)"
	jamKontrol := "08:00 - selesai"

	var val string
	err = s.db.QueryRowContext(ctx, `SELECT nilai FROM pengaturan WHERE kunci = 'nama_klinik'`).Scan(&val)
	if err == nil && val != "" {
		namaKlinik = val
	}
	err = s.db.QueryRowContext(ctx, `SELECT nilai FROM pengaturan WHERE kunci = 'jam_kontrol'`).Scan(&val)
	if err == nil && val != "" {
		jamKontrol = val
	}

	for _, jk := range schedules {
		tanggalStr := formatIndonesianDate(jk.TanggalKontrol)
		message := fmt.Sprintf(
			"Reminder Kontrol 🏥\n\nIbu %s, jangan lupa jadwal kontrol pada:\n\n📅 %s\n⏰ %s\n\nDi %s 😊\nTerima kasih.",
			jk.NamaPasien, tanggalStr, jamKontrol, namaKlinik,
		)

		judul := "Pengingat Jadwal Kontrol"
		jadwalID := jk.ID

		err := s.waGateway.SendMessage(jk.NoWaPasien, message)
		if err != nil {
			log.Printf(" Failed to send WA to %s: %v\n", jk.NoWaPasien, err)
			s.jadwalRepo.UpdateNotifStatus(ctx, jk.ID, "gagal")
			now := time.Now()
			s.notifRepo.Create(ctx, jk.PasienID, &jadwalID, judul, message, "gagal", &now)
		} else {
			log.Printf("WA sent to %s (%s)\n", jk.NamaPasien, jk.NoWaPasien)
			s.jadwalRepo.UpdateNotifStatus(ctx, jk.ID, "terkirim")
			now := time.Now()
			s.notifRepo.Create(ctx, jk.PasienID, &jadwalID, judul, message, "terkirim", &now)
		}
	}
}

func formatIndonesianDate(t time.Time) string {
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	months := []string{
		"",
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	dayName := days[t.Weekday()]
	dayNum := t.Day()
	monthName := months[t.Month()]
	year := t.Year()

	return fmt.Sprintf("%s, %d %s %d", dayName, dayNum, monthName, year)
}
