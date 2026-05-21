package models

import "time"

// Obat represents the obat (medicine master data) table.
type Obat struct {
	ID                int        `json:"id"`
	NamaObat          string     `json:"nama_obat"`
	Kategori          *string    `json:"kategori"`
	Satuan            string     `json:"satuan"`
	StokMinimum       int        `json:"stok_minimum"`
	JumlahStok        int        `json:"jumlah_stok"`
	TanggalKadaluarsa *time.Time `json:"tanggal_kadaluarsa"`
	BatasStokKritis   int        `json:"batas_stok_kritis"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
