package dto

// CreateObatRequest is the request body for adding a new medicine.
type CreateObatRequest struct {
	NamaObat    string `json:"nama_obat" validate:"required"`
	Kategori    string `json:"kategori" validate:"omitempty"`
	Satuan      string `json:"satuan" validate:"required"`
	StokMinimum int    `json:"stok_minimum" validate:"omitempty,gte=0"`
}

// UpdateObatRequest is the request body for updating a medicine.
type UpdateObatRequest struct {
	NamaObat          string `json:"nama_obat" validate:"omitempty"`
	Kategori          string `json:"kategori" validate:"omitempty"`
	Satuan            string `json:"satuan" validate:"omitempty"`
	StokMinimum       *int   `json:"stok_minimum" validate:"omitempty,gte=0"`
	BatasStokKritis   *int   `json:"batas_stok_kritis" validate:"omitempty,gte=0"`
	JumlahStok        *int   `json:"jumlah_stok" validate:"omitempty,gte=0"`
	TanggalKadaluarsa string `json:"tanggal_kadaluarsa" validate:"omitempty"`
}

// StokMasukRequest is the request body for adding incoming stock.
type StokMasukRequest struct {
	InventoriID       int    `json:"inventori_id" validate:"required"`
	Jumlah            int    `json:"jumlah" validate:"required,gt=0"`
	Keterangan        string `json:"keterangan" validate:"omitempty"`
	TanggalKadaluarsa string `json:"tanggal_kadaluarsa" validate:"omitempty"`
	BatchNumber       string `json:"batch_number" validate:"omitempty"`
}

// StokKeluarRequest is the request body for reducing stock (usage).
type StokKeluarRequest struct {
	InventoriID int    `json:"inventori_id" validate:"required"`
	Jumlah      int    `json:"jumlah" validate:"required,gt=0"`
	Keterangan  string `json:"keterangan" validate:"omitempty"`
}

// CreateObatWithInventoriRequest combines creating a medicine and its inventory entry.
type CreateObatWithInventoriRequest struct {
	NamaObat          string `json:"nama_obat" validate:"required"`
	Kategori          string `json:"kategori" validate:"omitempty"`
	Satuan            string `json:"satuan" validate:"required"`
	StokMinimum       int    `json:"stok_minimum" validate:"omitempty,gte=0"`
	BatasStokKritis   int    `json:"batas_stok_kritis" validate:"omitempty,gte=0"`
	JumlahStok        int    `json:"jumlah_stok" validate:"omitempty,gte=0"`
	TanggalKadaluarsa string `json:"tanggal_kadaluarsa" validate:"omitempty"`
	BatchNumber       string `json:"batch_number" validate:"omitempty"`
}
