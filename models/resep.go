package models

import "time"

// Resep represents the resep (prescription header) table.
type Resep struct {
	ID           int       `json:"id"`
	RekamMedisID int       `json:"rekam_medis_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Nested data
	Details []DetailResep `json:"details,omitempty"`
}
