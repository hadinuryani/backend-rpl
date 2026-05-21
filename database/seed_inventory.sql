-- ============================================
-- IC+ Klinik Bidan — Seed Data Obat & Inventori
-- Obat-obatan dunia kebidanan (real-world)
-- ============================================

-- Hapus data lama jika ada
DELETE FROM riwayat_stok;
DELETE FROM inventori;
DELETE FROM detail_resep;
DELETE FROM obat;

-- Reset auto increment
ALTER TABLE obat AUTO_INCREMENT = 1;
ALTER TABLE inventori AUTO_INCREMENT = 1;

-- ============================================
-- INSERT OBAT (Master Data Obat)
-- ============================================
INSERT INTO obat (nama_obat, kategori, satuan, stok_minimum) VALUES
-- === VITAMIN & SUPLEMEN KEHAMILAN ===
('Asam Folat 400 mcg', 'Vitamin & Suplemen', 'tablet', 50),
('Asam Folat 1 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Tablet Tambah Darah (Fe) 60 mg', 'Vitamin & Suplemen', 'tablet', 100),
('Ferrous Sulfate 300 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Ferrous Fumarate 200 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Kalsium Laktat 500 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Calcium Carbonate 500 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Vitamin B6 (Pyridoxine) 10 mg', 'Vitamin & Suplemen', 'tablet', 30),
('Vitamin B12 (Cyanocobalamin) 50 mcg', 'Vitamin & Suplemen', 'tablet', 30),
('Vitamin C 100 mg', 'Vitamin & Suplemen', 'tablet', 50),
('Vitamin C 500 mg', 'Vitamin & Suplemen', 'tablet', 30),
('Vitamin D3 1000 IU', 'Vitamin & Suplemen', 'tablet', 30),
('Vitamin K1 (Phytomenadione) 2 mg', 'Vitamin & Suplemen', 'tablet', 20),
('Vitamin K1 Injeksi 1 mg/0.5ml', 'Vitamin & Suplemen', 'ampul', 30),
('Multivitamin Ibu Hamil (Obimin AF)', 'Vitamin & Suplemen', 'kaplet', 50),
('DHA Omega-3 Fish Oil 500 mg', 'Vitamin & Suplemen', 'softgel', 30),
('Zinc 20 mg', 'Vitamin & Suplemen', 'tablet', 30),

-- === OBAT ANTI MUAL & MUNTAH (ANTIEMETIK) ===
('Ondansetron 4 mg', 'Antiemetik', 'tablet', 20),
('Ondansetron 8 mg', 'Antiemetik', 'tablet', 20),
('Domperidone 10 mg', 'Antiemetik', 'tablet', 20),
('Metoclopramide 10 mg', 'Antiemetik', 'tablet', 20),
('Dimenhidrinat 50 mg', 'Antiemetik', 'tablet', 20),
('Pyridoxine + Doxylamine (Bonjesta)', 'Antiemetik', 'tablet', 15),

-- === OBAT ANALGESIK & ANTIPIRETIK ===
('Paracetamol 500 mg', 'Analgesik & Antipiretik', 'tablet', 100),
('Paracetamol Sirup 120 mg/5ml', 'Analgesik & Antipiretik', 'botol', 20),
('Paracetamol Infus 1000 mg/100ml', 'Analgesik & Antipiretik', 'botol', 10),
('Ibuprofen 200 mg', 'Analgesik & Antipiretik', 'tablet', 30),
('Ibuprofen 400 mg', 'Analgesik & Antipiretik', 'tablet', 30),
('Asam Mefenamat 500 mg', 'Analgesik & Antipiretik', 'kaplet', 30),
('Metamizole (Antalgin) 500 mg', 'Analgesik & Antipiretik', 'tablet', 20),

-- === OBAT ANTIHIPERTENSI (PREEKLAMPSIA) ===
('Nifedipine 10 mg', 'Antihipertensi', 'tablet', 20),
('Nifedipine 30 mg SR', 'Antihipertensi', 'tablet', 15),
('Methyldopa 250 mg', 'Antihipertensi', 'tablet', 30),
('Methyldopa 500 mg', 'Antihipertensi', 'tablet', 20),
('MgSO4 (Magnesium Sulfat) 20% Injeksi', 'Antihipertensi', 'vial', 20),
('MgSO4 (Magnesium Sulfat) 40% Injeksi', 'Antihipertensi', 'vial', 15),
('Labetalol 100 mg', 'Antihipertensi', 'tablet', 15),

-- === OBAT UTEROTONIKA (KONTRAKSI RAHIM) ===
('Oksitosin Injeksi 10 IU/ml', 'Uterotonika', 'ampul', 30),
('Misoprostol 200 mcg', 'Uterotonika', 'tablet', 20),
('Methylergometrine 0.2 mg Tablet', 'Uterotonika', 'tablet', 20),
('Methylergometrine 0.2 mg Injeksi', 'Uterotonika', 'ampul', 20),
('Metergin (Methylergometrine) 0.125 mg', 'Uterotonika', 'tablet', 15),

-- === OBAT TOKOLITIK (ANTI KONTRAKSI PREMATUR) ===
('Nifedipine 10 mg (Tokolitik)', 'Tokolitik', 'tablet', 15),
('Isoxsuprine (Duvadilan) 10 mg', 'Tokolitik', 'tablet', 15),
('Terbutaline 2.5 mg', 'Tokolitik', 'tablet', 10),
('Terbutaline Injeksi 0.5 mg/ml', 'Tokolitik', 'ampul', 10),

-- === ANTIBIOTIK ===
('Amoxicillin 500 mg', 'Antibiotik', 'kapsul', 50),
('Amoxicillin Sirup 125 mg/5ml', 'Antibiotik', 'botol', 15),
('Ampicillin 500 mg', 'Antibiotik', 'kapsul', 30),
('Ampicillin Injeksi 1 g', 'Antibiotik', 'vial', 15),
('Cefadroxil 500 mg', 'Antibiotik', 'kapsul', 30),
('Cefixime 100 mg', 'Antibiotik', 'kapsul', 20),
('Cefixime 200 mg', 'Antibiotik', 'kapsul', 20),
('Ceftriaxone Injeksi 1 g', 'Antibiotik', 'vial', 10),
('Gentamicin Injeksi 80 mg/2ml', 'Antibiotik', 'ampul', 10),
('Eritromisin 500 mg', 'Antibiotik', 'tablet', 20),
('Azithromycin 500 mg', 'Antibiotik', 'tablet', 15),
('Metronidazole 500 mg', 'Antibiotik', 'tablet', 30),
('Metronidazole Infus 500 mg/100ml', 'Antibiotik', 'botol', 10),
('Clindamycin 300 mg', 'Antibiotik', 'kapsul', 20),
('Cotrimoxazole (TMP-SMX) 480 mg', 'Antibiotik', 'tablet', 20),
('Ciprofloxacin 500 mg', 'Antibiotik', 'tablet', 20),
('Doxycycline 100 mg', 'Antibiotik', 'kapsul', 15),

-- === ANTIJAMUR ===
('Fluconazole 150 mg', 'Antijamur', 'kapsul', 10),
('Nystatin Oral Drop 100.000 IU/ml', 'Antijamur', 'botol', 10),
('Miconazole Krim 2%', 'Antijamur', 'tube', 10),
('Clotrimazole Ovula 200 mg', 'Antijamur', 'ovula', 10),
('Ketoconazole 200 mg', 'Antijamur', 'tablet', 10),

-- === OBAT SALURAN CERNA ===
('Antasida DOEN (Al(OH)3 + Mg(OH)2)', 'Saluran Cerna', 'tablet', 30),
('Antasida Sirup', 'Saluran Cerna', 'botol', 15),
('Ranitidine 150 mg', 'Saluran Cerna', 'tablet', 20),
('Omeprazole 20 mg', 'Saluran Cerna', 'kapsul', 20),
('Lansoprazole 30 mg', 'Saluran Cerna', 'kapsul', 15),
('Sucralfate Sirup 500 mg/5ml', 'Saluran Cerna', 'botol', 10),
('Loperamide 2 mg', 'Saluran Cerna', 'kapsul', 15),
('Attapulgite 600 mg', 'Saluran Cerna', 'tablet', 20),
('Oralit (ORS)', 'Saluran Cerna', 'sachet', 50),
('Lactulose Sirup 3.35 g/5ml', 'Saluran Cerna', 'botol', 10),
('Bisacodyl 5 mg', 'Saluran Cerna', 'tablet', 15),

-- === OBAT ANTIALERGI & ANTIHISTAMIN ===
('Cetirizine 10 mg', 'Antihistamin', 'tablet', 20),
('Loratadine 10 mg', 'Antihistamin', 'tablet', 20),
('Chlorpheniramine Maleate (CTM) 4 mg', 'Antihistamin', 'tablet', 30),
('Dexamethasone 0.5 mg', 'Antihistamin', 'tablet', 20),
('Dexamethasone Injeksi 5 mg/ml', 'Antihistamin', 'ampul', 15),

-- === OBAT BATUK & SALURAN NAPAS ===
('Ambroxol 30 mg', 'Batuk & Saluran Napas', 'tablet', 30),
('Ambroxol Sirup 15 mg/5ml', 'Batuk & Saluran Napas', 'botol', 15),
('Bromhexine 8 mg', 'Batuk & Saluran Napas', 'tablet', 20),
('Guaifenesin 100 mg', 'Batuk & Saluran Napas', 'tablet', 20),
('Dextromethorphan 15 mg', 'Batuk & Saluran Napas', 'tablet', 15),
('Salbutamol 2 mg', 'Batuk & Saluran Napas', 'tablet', 15),
('Salbutamol Inhaler 100 mcg/puff', 'Batuk & Saluran Napas', 'inhaler', 5),

-- === OBAT ANTISEPTIK & TOPIKAL ===
('Povidone Iodine 10% (Betadine)', 'Antiseptik & Topikal', 'botol', 15),
('Chlorhexidine 0.5%', 'Antiseptik & Topikal', 'botol', 10),
('Alkohol 70%', 'Antiseptik & Topikal', 'botol', 20),
('Gentamicin Salep 0.1%', 'Antiseptik & Topikal', 'tube', 15),
('Silver Sulfadiazine Krim 1%', 'Antiseptik & Topikal', 'tube', 10),
('Kasa Steril 16x16 cm', 'Antiseptik & Topikal', 'pack', 30),
('Bioplacenton Gel', 'Antiseptik & Topikal', 'tube', 10),
('Sofra-tulle (Framycetin)', 'Antiseptik & Topikal', 'lembar', 20),
('Lidocaine Jelly 2%', 'Antiseptik & Topikal', 'tube', 10),

-- === CAIRAN INFUS ===
('NaCl 0.9% Infus 500 ml', 'Cairan Infus', 'botol', 20),
('Ringer Laktat (RL) 500 ml', 'Cairan Infus', 'botol', 20),
('Dextrose 5% Infus 500 ml', 'Cairan Infus', 'botol', 15),
('Dextrose 10% Infus 500 ml', 'Cairan Infus', 'botol', 10),
('Gelafundin 500 ml', 'Cairan Infus', 'botol', 5),

-- === ANESTESI LOKAL ===
('Lidocaine HCl 2% Injeksi', 'Anestesi Lokal', 'ampul', 20),
('Lidocaine + Epinefrin Injeksi', 'Anestesi Lokal', 'ampul', 15),
('Bupivacaine 0.5% Injeksi', 'Anestesi Lokal', 'vial', 10),

-- === OBAT EMERGENSI ===
('Epinefrin (Adrenalin) 1 mg/ml Injeksi', 'Obat Emergensi', 'ampul', 10),
('Atropin Sulfat 0.25 mg/ml Injeksi', 'Obat Emergensi', 'ampul', 10),
('Diazepam 5 mg/ml Injeksi', 'Obat Emergensi', 'ampul', 10),
('Dexamethasone 5 mg Injeksi (Emergensi)', 'Obat Emergensi', 'ampul', 10),
('Aminofilin Injeksi 24 mg/ml', 'Obat Emergensi', 'ampul', 5),
('Furosemide 20 mg/2ml Injeksi', 'Obat Emergensi', 'ampul', 10),

-- === OBAT KB (KONTRASEPSI) ===
('Pil KB Kombinasi (Levonorgestrel + EE)', 'Kontrasepsi', 'strip', 30),
('Pil KB Progestin Only (Minipil)', 'Kontrasepsi', 'strip', 20),
('Suntik KB 1 Bulan (Cyclofem)', 'Kontrasepsi', 'vial', 30),
('Suntik KB 3 Bulan (DMPA/Depo Provera)', 'Kontrasepsi', 'vial', 30),
('Implant KB (Implanon/Nexplanon)', 'Kontrasepsi', 'set', 10),
('IUD Copper T (AKDR)', 'Kontrasepsi', 'pcs', 10),
('Kondom', 'Kontrasepsi', 'pcs', 50),

-- === OBAT NEONATUS (BAYI BARU LAHIR) ===
('Salep Mata Chloramphenicol 1%', 'Obat Neonatus', 'tube', 15),
('Salep Mata Gentamicin 0.3%', 'Obat Neonatus', 'tube', 10),
('Hepatitis B Vaksin (HB-0)', 'Obat Neonatus', 'vial', 15),
('Vaksin BCG', 'Obat Neonatus', 'vial', 10),
('Vaksin Polio Oral (OPV)', 'Obat Neonatus', 'vial', 10),
('Gentian Violet 1%', 'Obat Neonatus', 'botol', 10),
('Minyak Telon', 'Obat Neonatus', 'botol', 15),
('Zalf Salep Tali Pusat', 'Obat Neonatus', 'tube', 15),

-- === OBAT LAKTASI (MENYUSUI) ===
('Domperidone 10 mg (Laktasi)', 'Obat Laktasi', 'tablet', 20),
('Lanolin Nipple Cream', 'Obat Laktasi', 'tube', 10),
('Laktafit (Suplemen ASI)', 'Obat Laktasi', 'kapsul', 15),

-- === OBAT ANTIANEMIA ===
('Asam Tranexamat 500 mg', 'Antianemia & Hemostatik', 'tablet', 20),
('Asam Tranexamat Injeksi 500 mg/5ml', 'Antianemia & Hemostatik', 'ampul', 10),
('Etamsylate 500 mg', 'Antianemia & Hemostatik', 'tablet', 15),
('Etamsylate Injeksi 250 mg/2ml', 'Antianemia & Hemostatik', 'ampul', 10),

-- === OBAT ANTI DIABETES GESTASIONAL ===
('Insulin Reguler (Actrapid) 100 IU/ml', 'Antidiabetes', 'vial', 5),
('Metformin 500 mg', 'Antidiabetes', 'tablet', 15),

-- === OBAT ANTI KEJANG ===
('MgSO4 40% (Anti Kejang Eklampsia)', 'Anti Kejang', 'vial', 10),
('Diazepam 5 mg Tablet', 'Anti Kejang', 'tablet', 15),
('Phenobarbital 30 mg', 'Anti Kejang', 'tablet', 10),

-- === ALAT HABIS PAKAI (BHP) ===
('Sarung Tangan Steril (Ukuran 7)', 'Alat Habis Pakai', 'pasang', 50),
('Sarung Tangan Steril (Ukuran 7.5)', 'Alat Habis Pakai', 'pasang', 50),
('Sarung Tangan Non-Steril', 'Alat Habis Pakai', 'box', 20),
('Spuit 1 cc', 'Alat Habis Pakai', 'pcs', 30),
('Spuit 3 cc', 'Alat Habis Pakai', 'pcs', 50),
('Spuit 5 cc', 'Alat Habis Pakai', 'pcs', 50),
('Spuit 10 cc', 'Alat Habis Pakai', 'pcs', 30),
('Infus Set Dewasa', 'Alat Habis Pakai', 'pcs', 15),
('IV Catheter No. 18', 'Alat Habis Pakai', 'pcs', 20),
('IV Catheter No. 20', 'Alat Habis Pakai', 'pcs', 20),
('IV Catheter No. 22', 'Alat Habis Pakai', 'pcs', 20),
('IV Catheter No. 24', 'Alat Habis Pakai', 'pcs', 20),
('Plester/Hansaplast', 'Alat Habis Pakai', 'roll', 15),
('Kapas Bulat Steril', 'Alat Habis Pakai', 'pack', 20),
('Catgut Chromic 2-0', 'Alat Habis Pakai', 'pack', 15),
('Catgut Chromic 3-0', 'Alat Habis Pakai', 'pack', 15),
('Benang Jahit Vicryl 2-0', 'Alat Habis Pakai', 'pack', 10),
('Benang Jahit Vicryl 3-0', 'Alat Habis Pakai', 'pack', 10),
('Kateter Foley No. 16', 'Alat Habis Pakai', 'pcs', 10),
('Kateter Foley No. 18', 'Alat Habis Pakai', 'pcs', 10),
('Urine Bag', 'Alat Habis Pakai', 'pcs', 15),
('Underpad Disposable', 'Alat Habis Pakai', 'pcs', 30),
('Masker Medis', 'Alat Habis Pakai', 'box', 10),
('Topi Operasi (Surgical Cap)', 'Alat Habis Pakai', 'pcs', 20),
('Apron Plastik Disposable', 'Alat Habis Pakai', 'pcs', 20);


-- ============================================
-- INSERT INVENTORI (Stok Inventory)
-- ============================================
-- Menggunakan subquery untuk ambil obat_id berdasarkan nama_obat

INSERT INTO inventori (obat_id, jumlah_stok, tanggal_kadaluarsa, batch_number, status_stok)
SELECT id, 
  CASE 
    WHEN kategori = 'Vitamin & Suplemen' THEN FLOOR(80 + RAND() * 120)
    WHEN kategori = 'Antibiotik' THEN FLOOR(40 + RAND() * 60)
    WHEN kategori = 'Alat Habis Pakai' THEN FLOOR(50 + RAND() * 150)
    WHEN kategori = 'Cairan Infus' THEN FLOOR(20 + RAND() * 30)
    WHEN kategori = 'Obat Emergensi' THEN FLOOR(10 + RAND() * 20)
    WHEN kategori = 'Kontrasepsi' THEN FLOOR(30 + RAND() * 50)
    WHEN kategori = 'Obat Neonatus' THEN FLOOR(15 + RAND() * 30)
    ELSE FLOOR(20 + RAND() * 80)
  END AS jumlah_stok,
  DATE_ADD(CURDATE(), INTERVAL FLOOR(6 + RAND() * 18) MONTH) AS tanggal_kadaluarsa,
  CONCAT('BTH-', DATE_FORMAT(CURDATE(), '%Y%m'), '-', LPAD(id, 4, '0')) AS batch_number,
  'aman' AS status_stok
FROM obat;

-- ============================================
-- UPDATE STATUS STOK berdasarkan jumlah & kadaluarsa
-- ============================================
UPDATE inventori i
JOIN obat o ON i.obat_id = o.id
SET i.status_stok = 
  CASE
    WHEN i.jumlah_stok = 0 THEN 'habis'
    WHEN i.jumlah_stok <= o.stok_minimum THEN 'hampir_habis'
    WHEN i.tanggal_kadaluarsa <= DATE_ADD(CURDATE(), INTERVAL 3 MONTH) THEN 'hampir_kadaluarsa'
    WHEN i.tanggal_kadaluarsa <= CURDATE() THEN 'kadaluarsa'
    ELSE 'aman'
  END;

-- ============================================
-- INSERT RIWAYAT STOK (Initial Stock Entry)
-- ============================================
INSERT INTO riwayat_stok (inventori_id, bidan_id, jenis_transaksi, jumlah, keterangan, created_at)
SELECT i.id, 1, 'masuk', i.jumlah_stok, 
  CONCAT('Stok awal - ', o.nama_obat, ' (Batch: ', i.batch_number, ')'),
  NOW()
FROM inventori i
JOIN obat o ON i.obat_id = o.id;

SELECT CONCAT('Total obat: ', COUNT(*)) AS info FROM obat
UNION ALL
SELECT CONCAT('Total inventori: ', COUNT(*)) FROM inventori
UNION ALL
SELECT CONCAT('Total riwayat stok: ', COUNT(*)) FROM riwayat_stok;
