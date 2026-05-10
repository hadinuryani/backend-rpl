package models

import "time"

// KlinikStatus represents the klinik_status table.
type KlinikStatus struct {
	ID        int       `json:"id"`
	BidanID   int       `json:"bidan_id"`
	Status    string    `json:"status"`
	Catatan   *string   `json:"catatan"`
	UpdatedAt time.Time `json:"updated_at"`
}
