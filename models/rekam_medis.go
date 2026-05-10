package models

import "time"

// RekamMedis represents the rekam_medis (medical records) table.
type RekamMedis struct {
	ID                int       `json:"id"`
	AntrianID         int       `json:"antrian_id"`
	BidanID           int       `json:"bidan_id"`
	KeluhanUtama      string    `json:"keluhan_utama"`
	TekananDarah      *string   `json:"tekanan_darah"`
	BeratBadan        *float64  `json:"berat_badan"`
	TinggiFundusUteri *float64  `json:"tinggi_fundus_uteri"`
	KondisiJanin      *string   `json:"kondisi_janin"`
	CatatanTambahan   *string   `json:"catatan_tambahan"`
	CreatedAt         time.Time `json:"created_at"`

	// Joined fields
	NamaPasien       string `json:"nama_pasien,omitempty"`
	TanggalKunjungan string `json:"tanggal_kunjungan,omitempty"`
	NamaBidan        string `json:"nama_bidan,omitempty"`
}
