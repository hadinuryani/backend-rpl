package dto

// CreateAntrianRequest is the request body for registering a new visit.
type CreateAntrianRequest struct {
	TanggalKunjungan string `json:"tanggal_kunjungan" validate:"required"`
	Keluhan          string `json:"keluhan" validate:"required,max=500"`
}
