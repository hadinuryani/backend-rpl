package models

import "time"

// Antrian represents the antrian (queue) table.
type Antrian struct {
	ID               int       `json:"id"`
	PasienID         int       `json:"pasien_id"`
	TanggalKunjungan time.Time `json:"tanggal_kunjungan"`
	NoAntrian        string    `json:"no_antrian"`
	Keluhan          string    `json:"keluhan"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Joined fields (not in DB, populated by queries)
	NamaPasien string `json:"nama_pasien,omitempty"`
	Umur       int    `json:"umur,omitempty"`
}
