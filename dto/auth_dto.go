package dto

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	NamaLengkap     string `json:"nama_lengkap" validate:"required,min=3"`
	TanggalLahir    string `json:"tanggal_lahir" validate:"required"`
	JenisKelamin    string `json:"jenis_kelamin" validate:"required,oneof=perempuan laki-laki"`
	Alamat          string `json:"alamat" validate:"required"`
	NoWa            string `json:"no_wa" validate:"required,min=10,max=15"`
	GolonganDarah   string `json:"golongan_darah" validate:"omitempty,oneof=A B AB O"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is the response body after successful login.
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// UserProfileResponse is the response for GET /auth/me.
type UserProfileResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`

	// Profile data (varies by role)
	Profile interface{} `json:"profile"`
}
