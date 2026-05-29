package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ic-plus-backend/models"
)

// MockUserRepo implements repositories.UserRepo for unit testing
type MockUserRepo struct {
	Users        map[int]*models.User
	UsersByEmail map[string]*models.User
	NextID       int
	OTPStore     map[int]struct {
		Code    string
		Expired time.Time
	}
	ErrorToReturn error
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		Users:        make(map[int]*models.User),
		UsersByEmail: make(map[string]*models.User),
		NextID:       1,
		OTPStore: make(map[int]struct {
			Code    string
			Expired time.Time
		}),
	}
}

func (m *MockUserRepo) Create(ctx context.Context, email, passwordHash, role string) (*models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	if _, ok := m.UsersByEmail[email]; ok {
		return nil, errors.New("email already exists")
	}
	user := &models.User{
		ID:           m.NextID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.Users[m.NextID] = user
	m.UsersByEmail[email] = user
	m.NextID++
	return user, nil
}

func (m *MockUserRepo) CreateTx(ctx context.Context, tx *sql.Tx, email, passwordHash, role string) (*models.User, error) {
	return m.Create(ctx, email, passwordHash, role)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.UsersByEmail[email], nil
}

func (m *MockUserRepo) FindByID(ctx context.Context, id int) (*models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.Users[id], nil
}

func (m *MockUserRepo) StoreOTP(ctx context.Context, userID int, otpCode string, expiredAt time.Time) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	m.OTPStore[userID] = struct {
		Code    string
		Expired time.Time
	}{Code: otpCode, Expired: expiredAt}
	return nil
}

func (m *MockUserRepo) FindByIDAndOTP(ctx context.Context, userID int, otpCode string) (*models.User, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	otp, ok := m.OTPStore[userID]
	if !ok || otp.Code != otpCode || otp.Expired.Before(time.Now()) {
		return nil, nil
	}
	return m.Users[userID], nil
}

func (m *MockUserRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if user, ok := m.Users[userID]; ok {
		user.PasswordHash = passwordHash
		user.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("user not found")
}

// MockPasienRepo implements repositories.PasienRepo
type MockPasienRepo struct {
	Pasiens       map[int]*models.Pasien
	PasiensByWa   map[string]*models.Pasien
	PasiensByUser map[int]*models.Pasien
	NextID        int
	ErrorToReturn error
}

func NewMockPasienRepo() *MockPasienRepo {
	return &MockPasienRepo{
		Pasiens:       make(map[int]*models.Pasien),
		PasiensByWa:   make(map[string]*models.Pasien),
		PasiensByUser: make(map[int]*models.Pasien),
		NextID:        1,
	}
}

func (m *MockPasienRepo) Create(ctx context.Context, userID int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) (*models.Pasien, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	pasien := &models.Pasien{
		ID:            m.NextID,
		UserID:        userID,
		NamaLengkap:   nama,
		TanggalLahir:  tglLahir,
		JenisKelamin:  jenisKelamin,
		Alamat:        &alamat,
		NoWa:          noWa,
		GolonganDarah: golDarah,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	m.Pasiens[m.NextID] = pasien
	m.PasiensByWa[noWa] = pasien
	m.PasiensByUser[userID] = pasien
	m.NextID++
	return pasien, nil
}

func (m *MockPasienRepo) CreateByBidan(ctx context.Context, tx *sql.Tx, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string, userID int) (*models.Pasien, error) {
	return m.Create(ctx, userID, nama, tglLahir, jenisKelamin, alamat, noWa, golDarah)
}

func (m *MockPasienRepo) FindByUserID(ctx context.Context, userID int) (*models.Pasien, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.PasiensByUser[userID], nil
}

func (m *MockPasienRepo) FindByID(ctx context.Context, id int) (*models.Pasien, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.Pasiens[id], nil
}

func (m *MockPasienRepo) FindByNoWa(ctx context.Context, noWa string) (*models.Pasien, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.PasiensByWa[noWa], nil
}

func (m *MockPasienRepo) Update(ctx context.Context, id int, nama string, tglLahir time.Time, jenisKelamin, alamat, noWa string, golDarah *string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	p, ok := m.Pasiens[id]
	if !ok {
		return errors.New("pasien not found")
	}
	p.NamaLengkap = nama
	p.TanggalLahir = tglLahir
	p.JenisKelamin = jenisKelamin
	p.Alamat = &alamat
	p.NoWa = noWa
	p.GolonganDarah = golDarah
	p.UpdatedAt = time.Now()
	return nil
}

func (m *MockPasienRepo) FindAll(ctx context.Context, search string, limit, offset int) ([]models.Pasien, int, error) {
	if m.ErrorToReturn != nil {
		return nil, 0, m.ErrorToReturn
	}
	var res []models.Pasien
	for _, p := range m.Pasiens {
		res = append(res, *p)
	}
	return res, len(res), nil
}

func (m *MockPasienRepo) UpdateIsHamil(ctx context.Context, id int, isHamil bool) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	p, ok := m.Pasiens[id]
	if !ok {
		return errors.New("pasien not found")
	}
	p.IsHamil = isHamil
	return nil
}

func (m *MockPasienRepo) UpdateIsHamilTx(ctx context.Context, tx *sql.Tx, id int, isHamil bool) error {
	return m.UpdateIsHamil(ctx, id, isHamil)
}

// MockAntrianRepo implements repositories.AntrianRepo
type MockAntrianRepo struct {
	CountTodayVal int
	CreateVal     *models.Antrian
	ErrorToReturn error
}

func (m *MockAntrianRepo) CountToday(ctx context.Context, dateStr string) (int, error) {
	return m.CountTodayVal, m.ErrorToReturn
}

func (m *MockAntrianRepo) Create(ctx context.Context, pasienID int, tanggal time.Time, noAntrian, keluhan string) (*models.Antrian, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	if m.CreateVal != nil {
		return m.CreateVal, nil
	}
	return &models.Antrian{
		ID:        1,
		PasienID:  pasienID,
		NoAntrian: noAntrian,
		Keluhan:   keluhan,
		Status:    "menunggu",
	}, nil
}

func (m *MockAntrianRepo) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Antrian, int, error) {
	return nil, 0, m.ErrorToReturn
}

func (m *MockAntrianRepo) FindByID(ctx context.Context, id int) (*models.Antrian, error) {
	return nil, m.ErrorToReturn
}

func (m *MockAntrianRepo) FindByDateAndStatus(ctx context.Context, dateStr string, status string, limit, offset int) ([]models.Antrian, int, error) {
	return nil, 0, m.ErrorToReturn
}

func (m *MockAntrianRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	return m.ErrorToReturn
}

func (m *MockAntrianRepo) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id int, status string) error {
	return m.ErrorToReturn
}

func (m *MockAntrianRepo) GetDashboardStats(ctx context.Context, dateStr string) (total, waiting, done int, err error) {
	return 0, 0, 0, m.ErrorToReturn
}

func (m *MockAntrianRepo) GetWeeklyVisitCounts(ctx context.Context, days int) ([]models.WeeklyVisit, error) {
	return nil, m.ErrorToReturn
}

// MockKlinikRepo implements repositories.KlinikRepo
type MockKlinikRepo struct {
	Status        string
	Catatan       string
	ErrorToReturn error
}

func (m *MockKlinikRepo) GetStatus(ctx context.Context) (status string, catatan string, err error) {
	return m.Status, m.Catatan, m.ErrorToReturn
}

func (m *MockKlinikRepo) SetStatus(ctx context.Context, bidanID int, status, catatan string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	m.Status = status
	m.Catatan = catatan
	return nil
}

func (m *MockKlinikRepo) GetStatusString(ctx context.Context) string {
	return m.Status
}

// MockJadwalRepo implements repositories.JadwalRepo
type MockJadwalRepo struct {
	Schedules     map[int]*models.JadwalKontrol
	NextID        int
	ErrorToReturn error
}

func NewMockJadwalRepo() *MockJadwalRepo {
	return &MockJadwalRepo{
		Schedules: make(map[int]*models.JadwalKontrol),
		NextID:    1,
	}
}

func (m *MockJadwalRepo) Create(ctx context.Context, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	jk := &models.JadwalKontrol{
		ID:             m.NextID,
		PasienID:       pasienID,
		BidanID:        bidanID,
		TanggalKontrol: tanggal,
		Catatan:        catatan,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	m.Schedules[m.NextID] = jk
	m.NextID++
	return jk, nil
}

func (m *MockJadwalRepo) CreateTx(ctx context.Context, tx *sql.Tx, pasienID, bidanID int, tanggal time.Time, catatan *string) (*models.JadwalKontrol, error) {
	return m.Create(ctx, pasienID, bidanID, tanggal, catatan)
}

func (m *MockJadwalRepo) FindAll(ctx context.Context, limit, offset int) ([]models.JadwalKontrol, int, error) {
	if m.ErrorToReturn != nil {
		return nil, 0, m.ErrorToReturn
	}
	var res []models.JadwalKontrol
	for _, s := range m.Schedules {
		res = append(res, *s)
	}
	return res, len(res), nil
}

func (m *MockJadwalRepo) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.JadwalKontrol, int, error) {
	if m.ErrorToReturn != nil {
		return nil, 0, m.ErrorToReturn
	}
	var res []models.JadwalKontrol
	for _, s := range m.Schedules {
		if s.PasienID == pasienID {
			res = append(res, *s)
		}
	}
	return res, len(res), nil
}

func (m *MockJadwalRepo) FindUpcomingForNotification(ctx context.Context, targetDate time.Time) ([]models.JadwalKontrol, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	var res []models.JadwalKontrol
	for _, s := range m.Schedules {
		// Compare date strings to simulate our fix
		if s.TanggalKontrol.Format("2006-01-02") == targetDate.Format("2006-01-02") && s.StatusNotifikasi == "belum" {
			res = append(res, *s)
		}
	}
	return res, nil
}

func (m *MockJadwalRepo) Update(ctx context.Context, id int, tanggal time.Time, catatan *string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if s, ok := m.Schedules[id]; ok {
		s.TanggalKontrol = tanggal
		s.Catatan = catatan
		s.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("schedule not found")
}

func (m *MockJadwalRepo) Delete(ctx context.Context, id int) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	delete(m.Schedules, id)
	return nil
}

func (m *MockJadwalRepo) UpdateNotifStatus(ctx context.Context, id int, status string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	if s, ok := m.Schedules[id]; ok {
		s.StatusNotifikasi = status
		s.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("schedule not found")
}

func (m *MockJadwalRepo) FindTodaySchedules(ctx context.Context, dateStr string) ([]models.JadwalKontrol, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	var res []models.JadwalKontrol
	for _, s := range m.Schedules {
		if s.TanggalKontrol.Format("2006-01-02") == dateStr {
			res = append(res, *s)
		}
	}
	return res, nil
}

// MockNotifikasiRepo implements repositories.NotifikasiRepo
type MockNotifikasiRepo struct {
	Notifs        []*models.Notifikasi
	ErrorToReturn error
}

func NewMockNotifikasiRepo() *MockNotifikasiRepo {
	return &MockNotifikasiRepo{
		Notifs: make([]*models.Notifikasi, 0),
	}
}

func (m *MockNotifikasiRepo) Create(ctx context.Context, pasienID int, jadwalKontrolID *int, judul, pesan, statusKirim string, sentAt *time.Time) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	m.Notifs = append(m.Notifs, &models.Notifikasi{
		ID:              len(m.Notifs) + 1,
		PasienID:        pasienID,
		JadwalKontrolID: jadwalKontrolID,
		Judul:           judul,
		Pesan:           pesan,
		StatusKirim:     statusKirim,
		SentAt:          sentAt,
	})
	return nil
}

func (m *MockNotifikasiRepo) FindByPasienID(ctx context.Context, pasienID, limit, offset int) ([]models.Notifikasi, int, error) {
	return nil, 0, m.ErrorToReturn
}

func (m *MockNotifikasiRepo) MarkAsRead(ctx context.Context, id, pasienID int) error {
	return m.ErrorToReturn
}

// MockPengaturanRepo implements repositories.PengaturanRepo
type MockPengaturanRepo struct {
	Settings      map[string]string
	ErrorToReturn error
}

func NewMockPengaturanRepo() *MockPengaturanRepo {
	return &MockPengaturanRepo{
		Settings: make(map[string]string),
	}
}

func (m *MockPengaturanRepo) GetAll(ctx context.Context) (map[string]string, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}
	return m.Settings, nil
}

func (m *MockPengaturanRepo) Upsert(ctx context.Context, kunci, nilai string) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	m.Settings[kunci] = nilai
	return nil
}

func (m *MockPengaturanRepo) GetByKey(ctx context.Context, key string) (string, error) {
	if m.ErrorToReturn != nil {
		return "", m.ErrorToReturn
	}
	return m.Settings[key], nil
}
