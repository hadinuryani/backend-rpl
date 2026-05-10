package services

import (
	"context"
	"fmt"
	"log"
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
}

func NewSchedulerService(jadwalRepo *repositories.JadwalRepository, notifRepo *repositories.NotifikasiRepository, waGateway WAGateway) *SchedulerService {
	// Use WIB timezone (UTC+7)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	c := cron.New(cron.WithLocation(loc))
	return &SchedulerService{
		jadwalRepo: jadwalRepo,
		notifRepo:  notifRepo,
		waGateway:  waGateway,
		cron:       c,
	}
}

// Start begins the cron scheduler with the H-1 notification job.
func (s *SchedulerService) Start() {
	// Run every day at 08:00 WIB
	s.cron.AddFunc("0 8 * * *", func() {
		log.Println("Running H-1 notification job...")
		s.runH1NotificationJob()
	})
	s.cron.Start()
	log.Println(" Scheduler started (H-1 notifications at 08:00 WIB)")
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

	for _, jk := range schedules {
		tanggalStr := jk.TanggalKontrol.Format("2 January 2006")
		message := fmt.Sprintf(
			"Halo %s 👋\nIni adalah pengingat dari Klinik Indah Care Plus (IC+).\nBesok, %s, Anda memiliki jadwal kontrol.\nHarap segera lakukan pendaftaran kunjungan melalui sistem kami.\nSampai jumpa! 🌿",
			jk.NamaPasien, tanggalStr,
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
