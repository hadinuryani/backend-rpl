-- ============================================
-- IC+ Klinik Bidan — Database Schema (MySQL)
-- ============================================

-- 1. Users (authentication)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('pasien', 'bidan')),
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 2. Pasien (patient profile)
CREATE TABLE pasien (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    nama_lengkap VARCHAR(255) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin VARCHAR(20) NOT NULL CHECK (jenis_kelamin IN ('perempuan', 'laki-laki')),
    alamat TEXT,
    no_wa VARCHAR(20) NOT NULL,
    golongan_darah VARCHAR(5),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Bidan (midwife profile)
CREATE TABLE bidan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    nama_lengkap VARCHAR(255) NOT NULL,
    no_str VARCHAR(100),
    no_wa VARCHAR(20),
    foto_profil VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 4. Klinik Status
CREATE TABLE klinik_status (
    id INT AUTO_INCREMENT PRIMARY KEY,
    bidan_id INT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('buka', 'tutup', 'libur')),
    catatan TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (bidan_id) REFERENCES bidan(id)
);

-- 5. Antrian (queue)
CREATE TABLE antrian (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pasien_id INT NOT NULL,
    tanggal_kunjungan DATE NOT NULL,
    no_antrian VARCHAR(10) NOT NULL,
    keluhan TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'menunggu' CHECK (status IN ('menunggu', 'selesai', 'batal')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (pasien_id) REFERENCES pasien(id)
);

-- 6. Rekam Medis (medical records)
CREATE TABLE rekam_medis (
    id INT AUTO_INCREMENT PRIMARY KEY,
    antrian_id INT UNIQUE NOT NULL,
    bidan_id INT NOT NULL,
    keluhan_utama TEXT NOT NULL,
    tekanan_darah VARCHAR(20),
    berat_badan DECIMAL(5,2),
    tinggi_fundus_uteri DECIMAL(5,2),
    kondisi_janin TEXT,
    catatan_tambahan TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (antrian_id) REFERENCES antrian(id),
    FOREIGN KEY (bidan_id) REFERENCES bidan(id)
);

-- 7. Resep (prescription header)
CREATE TABLE resep (
    id INT AUTO_INCREMENT PRIMARY KEY,
    rekam_medis_id INT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rekam_medis_id) REFERENCES rekam_medis(id)
);

-- 8. Obat (medicine master data)
CREATE TABLE obat (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama_obat VARCHAR(255) NOT NULL,
    kategori VARCHAR(100),
    satuan VARCHAR(50) NOT NULL,
    stok_minimum INT DEFAULT 10,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 9. Detail Resep (prescription items)
CREATE TABLE detail_resep (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resep_id INT NOT NULL,
    obat_id INT NOT NULL,
    dosis VARCHAR(100) NOT NULL,
    aturan_pakai VARCHAR(255) NOT NULL,
    catatan TEXT,
    FOREIGN KEY (resep_id) REFERENCES resep(id) ON DELETE CASCADE,
    FOREIGN KEY (obat_id) REFERENCES obat(id)
);

-- 10. Jadwal Kontrol (control schedule)
CREATE TABLE jadwal_kontrol (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pasien_id INT NOT NULL,
    bidan_id INT NOT NULL,
    tanggal_kontrol DATE NOT NULL,
    catatan TEXT,
    status_notifikasi VARCHAR(20) DEFAULT 'belum' CHECK (status_notifikasi IN ('belum', 'terkirim', 'gagal')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (pasien_id) REFERENCES pasien(id),
    FOREIGN KEY (bidan_id) REFERENCES bidan(id)
);

-- 11. Notifikasi (notifications)
CREATE TABLE notifikasi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pasien_id INT NOT NULL,
    jadwal_kontrol_id INT NULL,
    judul VARCHAR(255) NOT NULL,
    pesan TEXT NOT NULL,
    channel VARCHAR(50) DEFAULT 'whatsapp',
    status_kirim VARCHAR(20) DEFAULT 'terkirim' CHECK (status_kirim IN ('terkirim', 'gagal', 'pending')),
    is_read TINYINT(1) DEFAULT 0,
    sent_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pasien_id) REFERENCES pasien(id),
    FOREIGN KEY (jadwal_kontrol_id) REFERENCES jadwal_kontrol(id)
);

-- 12. Inventori (inventory)
CREATE TABLE inventori (
    id INT AUTO_INCREMENT PRIMARY KEY,
    obat_id INT UNIQUE NOT NULL,
    jumlah_stok INT NOT NULL DEFAULT 0,
    tanggal_kadaluarsa DATE,
    batch_number VARCHAR(100),
    status_stok VARCHAR(30) DEFAULT 'aman'
        CHECK (status_stok IN ('aman', 'hampir_habis', 'habis', 'hampir_kadaluarsa', 'kadaluarsa')),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (obat_id) REFERENCES obat(id)
);

-- 13. Riwayat Stok (stock history)
CREATE TABLE riwayat_stok (
    id INT AUTO_INCREMENT PRIMARY KEY,
    inventori_id INT NOT NULL,
    bidan_id INT NOT NULL,
    jenis_transaksi VARCHAR(10) NOT NULL CHECK (jenis_transaksi IN ('masuk', 'keluar')),
    jumlah INT NOT NULL,
    keterangan TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventori_id) REFERENCES inventori(id),
    FOREIGN KEY (bidan_id) REFERENCES bidan(id)
);

-- ============================================
-- INDEXES
-- ============================================
CREATE INDEX idx_antrian_pasien_id ON antrian(pasien_id);
CREATE INDEX idx_antrian_tanggal ON antrian(tanggal_kunjungan);
CREATE INDEX idx_antrian_status ON antrian(status);
CREATE INDEX idx_jadwal_tanggal ON jadwal_kontrol(tanggal_kontrol);
CREATE INDEX idx_jadwal_status_notif ON jadwal_kontrol(status_notifikasi);
CREATE INDEX idx_notifikasi_pasien ON notifikasi(pasien_id);
CREATE INDEX idx_rekam_medis_antrian ON rekam_medis(antrian_id);
