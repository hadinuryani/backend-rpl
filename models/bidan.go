package models

import "time"

// Bidan represents the bidan (midwife) table.
type Bidan struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	NamaLengkap string   `json:"nama_lengkap"`
	NoSTR       *string   `json:"no_str"`
	NoWa        *string   `json:"no_wa"`
	FotoProfil  *string   `json:"foto_profil"`
	CreatedAt   time.Time `json:"created_at"`
}
