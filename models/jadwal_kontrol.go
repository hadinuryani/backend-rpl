package models

import "time"

// JadwalKontrol represents the jadwal_kontrol (control schedule) table.
type JadwalKontrol struct {
	ID               int       `json:"id"`
	PasienID         int       `json:"pasien_id"`
	BidanID          int       `json:"bidan_id"`
	TanggalKontrol   time.Time `json:"tanggal_kontrol"`
	Catatan          *string   `json:"catatan"`
	StatusNotifikasi string    `json:"status_notifikasi"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Joined fields
	NamaPasien string `json:"nama_pasien,omitempty"`
	NoWaPasien string `json:"no_wa_pasien,omitempty"`
}
