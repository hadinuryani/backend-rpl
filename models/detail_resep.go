package models

// DetailResep represents the detail_resep (prescription line items) table.
type DetailResep struct {
	ID          int     `json:"id"`
	ResepID     int     `json:"resep_id"`
	ObatID      int     `json:"obat_id"`
	Jumlah      int     `json:"jumlah"`
	Dosis       string  `json:"dosis"`
	AturanPakai string  `json:"aturan_pakai"`
	Catatan     *string `json:"catatan"`

	// Joined fields
	NamaObat string `json:"nama_obat,omitempty"`
}
