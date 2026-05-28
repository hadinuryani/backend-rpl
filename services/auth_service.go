package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repositories.UserRepo
	pasienRepo repositories.PasienRepo
	waGateway  WAGateway
}

func NewAuthService(userRepo repositories.UserRepo, pasienRepo repositories.PasienRepo, waGateway WAGateway) *AuthService {
	return &AuthService{userRepo: userRepo, pasienRepo: pasienRepo, waGateway: waGateway}
}

func (s *AuthService) Register(ctx context.Context, email, password, nama, tglLahirStr, jenisKelamin, alamat, noWa, golDarah string) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("cek email: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, email, string(hash), "pasien")
	if err != nil {
		return nil, err
	}

	tglLahir, err := time.Parse("2006-01-02", tglLahirStr)
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
	}

	var golDarahPtr *string
	if golDarah != "" {
		golDarahPtr = &golDarah
	}

	_, err = s.pasienRepo.Create(ctx, user.ID, nama, tglLahir, jenisKelamin, alamat, noWa, golDarahPtr)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, fmt.Errorf("email atau password salah")
	}
	if !user.IsActive {
		return "", nil, fmt.Errorf("akun tidak aktif")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, fmt.Errorf("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID int) (*models.User, interface{}, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, nil, fmt.Errorf("user tidak ditemukan")
	}

	if user.Role == "pasien" {
		pasien, err := s.pasienRepo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, nil, err
		}
		return user, pasien, nil
	}

	return user, nil, nil
}

func generateOTP() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 || err != nil {
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (s *AuthService) ForgotPassword(ctx context.Context, noWa string) error {
	pasien, err := s.pasienRepo.FindByNoWa(ctx, noWa)
	if err != nil {
		return err
	}
	if pasien == nil {
		return fmt.Errorf("nomor WhatsApp tidak terdaftar")
	}

	user, err := s.userRepo.FindByID(ctx, pasien.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("akun pengguna tidak ditemukan")
	}

	otp := generateOTP()
	expiredAt := time.Now().Add(15 * time.Minute)

	err = s.userRepo.StoreOTP(ctx, user.ID, otp, expiredAt)
	if err != nil {
		return fmt.Errorf("gagal menyimpan OTP: %w", err)
	}

	message := fmt.Sprintf(
		"Halo 👋\nIni adalah kode OTP untuk melakukan reset password akun Klinik Indah Care Plus (IC+) Anda:\n\n*🔑 %s*\n\nKode ini berlaku selama 15 menit. Jangan bagikan kode ini kepada siapa pun.",
		otp,
	)

	err = s.waGateway.SendMessage(noWa, message)
	if err != nil {
		return fmt.Errorf("gagal mengirimkan OTP ke nomor %s: %w", noWa, err)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, noWa, otpCode, newPassword string) error {
	pasien, err := s.pasienRepo.FindByNoWa(ctx, noWa)
	if err != nil {
		return err
	}
	if pasien == nil {
		return fmt.Errorf("nomor WhatsApp tidak terdaftar")
	}

	user, err := s.userRepo.FindByIDAndOTP(ctx, pasien.UserID, otpCode)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("kode OTP salah atau telah kedaluwarsa")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("gagal memproses password baru: %w", err)
	}

	err = s.userRepo.UpdatePassword(ctx, user.ID, string(hash))
	if err != nil {
		return fmt.Errorf("gagal memperbarui password: %w", err)
	}

	return nil
}
