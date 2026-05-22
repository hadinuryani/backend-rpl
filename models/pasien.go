package models

import "time"

// Pasien represents the pasien (patient) table.
type Pasien struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	NamaLengkap    string    `json:"nama_lengkap"`
	TanggalLahir   time.Time `json:"tanggal_lahir"`
	JenisKelamin   string    `json:"jenis_kelamin"`
	Alamat         *string   `json:"alamat"`
	NoWa           string    `json:"no_wa"`
	GolonganDarah  *string   `json:"golongan_darah"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Umur           int       `json:"umur"`
}
