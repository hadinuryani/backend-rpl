package dto

// CreateJadwalRequest is the request body for creating a control schedule.
type CreateJadwalRequest struct {
	PasienID       int    `json:"pasien_id" validate:"required"`
	TanggalKontrol string `json:"tanggal_kontrol" validate:"required"`
	Catatan        string `json:"catatan" validate:"omitempty"`
}

// UpdateJadwalRequest is the request body for updating a control schedule.
type UpdateJadwalRequest struct {
	TanggalKontrol string `json:"tanggal_kontrol" validate:"omitempty"`
	Catatan        string `json:"catatan" validate:"omitempty"`
}
