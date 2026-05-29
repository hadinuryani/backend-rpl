package dto

// CreateRekamMedisRequest is the request body for creating a medical record + prescription.
// When resep items are provided, the antrian is auto-completed in a single transaction.
type CreateRekamMedisRequest struct {
	AntrianID         int              `json:"antrian_id" validate:"required"`
	KeluhanUtama      string           `json:"keluhan_utama" validate:"required"`
	TekananDarah      string           `json:"tekanan_darah" validate:"omitempty"`
	BeratBadan        *float64         `json:"berat_badan" validate:"omitempty,gt=0"`
	TinggiFundusUteri *float64         `json:"tinggi_fundus_uteri" validate:"omitempty"`
	KondisiJanin      string           `json:"kondisi_janin" validate:"omitempty"`
	CatatanTambahan   string           `json:"catatan_tambahan" validate:"omitempty"`
	Resep             []DetailResepDTO `json:"resep" validate:"required,min=1,dive"`
	PerluKontrol      bool             `json:"perlu_kontrol"`
	TanggalKontrol    string           `json:"tanggal_kontrol" validate:"omitempty"`
	CatatanKontrol    string           `json:"catatan_kontrol" validate:"omitempty"`
	IsHamil           *bool            `json:"is_hamil" validate:"omitempty"`
}

// DetailResepDTO represents a single prescription line item in the request.
type DetailResepDTO struct {
	ObatID      int    `json:"obat_id" validate:"required"`
	Jumlah      int    `json:"jumlah" validate:"required,gte=1"`
	Dosis       string `json:"dosis" validate:"required"`
	AturanPakai string `json:"aturan_pakai" validate:"required"`
	Catatan     string `json:"catatan" validate:"omitempty"`
}

// UpdateRekamMedisRequest is the request body for updating a medical record.
type UpdateRekamMedisRequest struct {
	KeluhanUtama      string   `json:"keluhan_utama" validate:"omitempty"`
	TekananDarah      string   `json:"tekanan_darah" validate:"omitempty"`
	BeratBadan        *float64 `json:"berat_badan" validate:"omitempty,gt=0"`
	TinggiFundusUteri *float64 `json:"tinggi_fundus_uteri" validate:"omitempty"`
	KondisiJanin      string   `json:"kondisi_janin" validate:"omitempty"`
	CatatanTambahan   string   `json:"catatan_tambahan" validate:"omitempty"`
}
