package dto

// UpdatePasienRequest is the request body for updating patient profile.
type UpdatePasienRequest struct {
	NamaLengkap   string `json:"nama_lengkap" validate:"omitempty,min=3"`
	TanggalLahir  string `json:"tanggal_lahir" validate:"omitempty"`
	JenisKelamin  string `json:"jenis_kelamin" validate:"omitempty,oneof=perempuan laki-laki"`
	Alamat        string `json:"alamat" validate:"omitempty"`
	NoWa          string `json:"no_wa" validate:"omitempty,min=10,max=15"`
	GolonganDarah string `json:"golongan_darah" validate:"omitempty,oneof=A B AB O"`
}

// CreatePasienByBidanRequest is the request body for bidan manually adding a patient.
type CreatePasienByBidanRequest struct {
	NamaLengkap   string `json:"nama_lengkap" validate:"required,min=3"`
	TanggalLahir  string `json:"tanggal_lahir" validate:"required"`
	JenisKelamin  string `json:"jenis_kelamin" validate:"required,oneof=perempuan laki-laki"`
	Alamat        string `json:"alamat" validate:"omitempty"`
	NoWa          string `json:"no_wa" validate:"required,min=10,max=15"`
	GolonganDarah string `json:"golongan_darah" validate:"omitempty,oneof=A B AB O"`
}
