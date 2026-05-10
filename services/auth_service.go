package services

import (
	"context"
	"fmt"
	"time"

	"ic-plus-backend/models"
	"ic-plus-backend/repositories"
	"ic-plus-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	pasienRepo *repositories.PasienRepository
}

func NewAuthService(userRepo *repositories.UserRepository, pasienRepo *repositories.PasienRepository) *AuthService {
	return &AuthService{userRepo: userRepo, pasienRepo: pasienRepo}
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
