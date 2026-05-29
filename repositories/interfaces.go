package repositories

import (
	"context"
	"database/sql"
	"time"

	"ic-plus-backend/models"
)

// UserRepo defines the contract for user data access.
type UserRepo interface {
	Create(ctx context.Context, email, passwordHash, role string) (*models.User, error)
	CreateTx(ctx context.Context, tx *sql.Tx, email, passwordHash, role string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	StoreOTP(ctx context.Context, userID int, otpCode string, expiredAt time.Time) error
	FindByIDAndOTP(ctx context.Context, userID int, otpCode string) (*models.User, error)
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

// PasienRepo defines the contract for patient data access.
type PasienRepo interface {
	Create(ctx context.Context, userID int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) (*models.Pasien, error)
	CreateByBidan(ctx context.Context, tx *sql.Tx, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string, userID int) (*models.Pasien, error)
	FindByUserID(ctx context.Context, userID int) (*models.Pasien, error)
	FindByID(ctx context.Context, id int) (*models.Pasien, error)
	FindByNoWa(ctx context.Context, noWa string) (*models.Pasien, error)
	Update(ctx context.Context, id int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) error
	FindAll(ctx context.Context, search string, limit, offset int) ([]models.Pasien, int, error)
	UpdateIsHamil(ctx context.Context, id int, isHamil bool) error
	UpdateIsHamilTx(ctx context.Context, tx *sql.Tx, id int, isHamil bool) error
}

// AntrianRepo defines the contract for queue data access.
type AntrianRepo interface {
	CountToday(ctx context.Context, dateStr string) (int, error)
	Create(ctx context.Context, pasienID int, tanggal time.Time, noAntrian, keluhan string) (*models.Antrian, error)
	FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Antrian, int, error)
	FindByID(ctx context.Context, id int) (*models.Antrian, error)
	FindByDateAndStatus(ctx context.Context, dateStr string, status string, limit, offset int) ([]models.Antrian, int, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int, status string) error
	GetDashboardStats(ctx context.Context, dateStr string) (total, waiting, done int, err error)
	GetWeeklyVisitCounts(ctx context.Context, days int) ([]models.WeeklyVisit, error)
}

// RekamMedisRepo defines the contract for medical record data access.
type RekamMedisRepo interface {
	CreateWithResepTx(ctx context.Context, tx *sql.Tx, bidanID, antrianID int, keluhanUtama string, tekananDarah *string, beratBadan, tinggiFundus *float64, kondisiJanin, catatan *string, details []struct {
		ObatID, Jumlah  int
		Dosis, AturanPakai string
		Catatan            *string
	}) (int, int, error)
	FindByID(ctx context.Context, id int) (*models.RekamMedis, error)
	FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.RekamMedis, int, error)
	Update(ctx context.Context, id int, keluhanUtama string, tekananDarah *string, beratBadan, tinggiFundus *float64, kondisiJanin, catatan *string) error
	FindResepByRekamMedisID(ctx context.Context, rmID int) (*models.Resep, error)
}

// JadwalRepo defines the contract for control schedule data access.
type JadwalRepo interface {
	Create(ctx context.Context, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error)
	CreateTx(ctx context.Context, tx *sql.Tx, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.JadwalKontrol, int, error)
	FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.JadwalKontrol, int, error)
	FindUpcomingForNotification(ctx context.Context, targetDate time.Time) ([]models.JadwalKontrol, error)
	Update(ctx context.Context, id int, tanggal time.Time, catatan *string) error
	Delete(ctx context.Context, id int) error
	UpdateNotifStatus(ctx context.Context, id int, status string) error
	FindTodaySchedules(ctx context.Context, dateStr string) ([]models.JadwalKontrol, error)
}

// NotifikasiRepo defines the contract for notification data access.
type NotifikasiRepo interface {
	Create(ctx context.Context, pasienID int, jadwalKontrolID *int, judul, pesan, statusKirim string, sentAt *time.Time) error
	FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Notifikasi, int, error)
	MarkAsRead(ctx context.Context, id, pasienID int) error
}

// InventoriRepo defines the contract for inventory data access.
type InventoriRepo interface {
	FindAllObat(ctx context.Context, search string, status string, limit, offset int) ([]models.Obat, int, error)
	CreateObatWithInventori(ctx context.Context, nama, kategori, satuan string, stokMin, jumlahStok int, tanggalKadaluarsa *time.Time, batchNumber *string) (*models.Obat, error)
	UpdateObat(ctx context.Context, id int, nama, kategori, satuan string, stokMin *int, jumlahStok *int, tglKadaluarsa *time.Time) error
	DeleteObat(ctx context.Context, id int) error
	FindAllInventori(ctx context.Context, statusFilter string, limit, offset int) ([]models.Inventori, int, error)
	StokMasuk(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan *string, tanggalKadaluarsa *time.Time, batchNumber *string) error
	StokKeluar(ctx context.Context, inventoriID, bidanID, jumlah int, keterangan *string) error
	GetRiwayat(ctx context.Context, limit, offset int) ([]models.RiwayatStok, int, error)
	CountCritical(ctx context.Context) (int, error)
	GetCriticalMedicines(ctx context.Context) ([]models.Inventori, error)
	GetStockStatusSummary(ctx context.Context) (map[string]int, error)
}

// KlinikRepo defines the contract for clinic status data access.
type KlinikRepo interface {
	GetStatus(ctx context.Context) (status string, catatan string, err error)
	SetStatus(ctx context.Context, bidanID int, status, catatan string) error
	GetStatusString(ctx context.Context) string
}

// BidanRepo defines the contract for bidan (midwife) data access.
type BidanRepo interface {
	FindBidanIDByUserID(ctx context.Context, userID int) (int, error)
}

// PengaturanRepo defines the contract for application settings data access.
type PengaturanRepo interface {
	GetAll(ctx context.Context) (map[string]string, error)
	Upsert(ctx context.Context, kunci, nilai string) error
	GetByKey(ctx context.Context, kunci string) (string, error)
}

// MonitorRepo defines the contract for visit monitoring data access.
type MonitorRepo interface {
	GetVisitHistory(ctx context.Context, from, to time.Time, search string, limit, offset int) ([]VisitRecord, int, error)
}

// VisitRecord is a data struct for monitor visit history results.
type VisitRecord struct {
	ID           int    `json:"id"`
	PasienID     int    `json:"pasien_id"`
	NoAntrian    string `json:"no_antrian"`
	NamaPasien   string `json:"nama_pasien"`
	TanggalDaftar string `json:"tanggal_daftar"`
	Keluhan      string `json:"keluhan"`
	Status       string `json:"status"`
	RekamMedisID *int   `json:"rekam_medis_id"`
}
