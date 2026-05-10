package models

import "time"

// Inventori represents the inventori (inventory) table.
type Inventori struct {
	ID                int        `json:"id"`
	ObatID            int        `json:"obat_id"`
	JumlahStok        int        `json:"jumlah_stok"`
	TanggalKadaluarsa *time.Time `json:"tanggal_kadaluarsa"`
	BatchNumber       *string    `json:"batch_number"`
	StatusStok        string     `json:"status_stok"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Joined fields
	NamaObat    string  `json:"nama_obat,omitempty"`
	Kategori    *string `json:"kategori,omitempty"`
	Satuan      string  `json:"satuan,omitempty"`
	StokMinimum int     `json:"stok_minimum,omitempty"`
}

// RiwayatStok represents the riwayat_stok (stock history) table.
type RiwayatStok struct {
	ID              int       `json:"id"`
	InventoriID     int       `json:"inventori_id"`
	BidanID         int       `json:"bidan_id"`
	JenisTransaksi  string    `json:"jenis_transaksi"`
	Jumlah          int       `json:"jumlah"`
	Keterangan      *string   `json:"keterangan"`
	CreatedAt       time.Time `json:"created_at"`

	// Joined fields
	NamaObat  string `json:"nama_obat,omitempty"`
	NamaBidan string `json:"nama_bidan,omitempty"`
}
