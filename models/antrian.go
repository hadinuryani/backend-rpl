package models

import (
	"encoding/json"
	"time"
)

// DateOnly is a date-only type that serializes to "YYYY-MM-DD" in JSON.
type DateOnly struct {
	time.Time
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(d.Format("2006-01-02"))
}

// Antrian represents the antrian (queue) table.
type Antrian struct {
	ID               int       `json:"id"`
	PasienID         int       `json:"pasien_id"`
	TanggalKunjungan DateOnly  `json:"tanggal_kunjungan"`
	NoAntrian        string    `json:"no_antrian"`
	Keluhan          string    `json:"keluhan"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Joined fields (not in DB, populated by queries)
	NamaPasien    string `json:"nama_pasien,omitempty"`
	Umur          int    `json:"umur,omitempty"`
	GolonganDarah string `json:"golongan_darah,omitempty"`
	Alamat        string `json:"alamat,omitempty"`
	NoWa          string `json:"no_wa,omitempty"`
	JenisKelamin  string `json:"jenis_kelamin,omitempty"`
	IsHamil       bool   `json:"is_hamil,omitempty"`
}

// WeeklyVisit represents visit count for a day.
type WeeklyVisit struct {
	Tanggal string `json:"tanggal"`
	Hari    string `json:"hari"`
	Jumlah  int    `json:"jumlah"`
}

