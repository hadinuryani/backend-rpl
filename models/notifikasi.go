package models

import "time"

// Notifikasi represents the notifikasi (notifications) table.
type Notifikasi struct {
	ID              int        `json:"id"`
	PasienID        int        `json:"pasien_id"`
	JadwalKontrolID *int       `json:"jadwal_kontrol_id"`
	Judul           string     `json:"judul"`
	Pesan           string     `json:"pesan"`
	Channel         string     `json:"channel"`
	StatusKirim     string     `json:"status_kirim"`
	IsRead          bool       `json:"is_read"`
	SentAt          *time.Time `json:"sent_at"`
	CreatedAt       time.Time  `json:"created_at"`
}
