package config

// JenisKelamin - Jenis kelamin
var JenisKelamin = map[int]string{
	1: "LAKI-LAKI",
	2: "PEREMPUAN",
	3: "LAINNYA",
}

// JK - Singkatan jenis kelamin
var JK = map[int]string{
	1: "Lk",
	2: "Pr",
}

// Agama - Daftar agama
var Agama = map[int]string{
	1: "ISLAM",
	2: "KRISTEN",
	3: "KATHOLIK",
	4: "HINDU",
	5: "BUDHA",
	6: "KHONGHUCU",
	7: "KEPERCAYAAN TERHADAP TUHAN YME / LAINNYA",
}

// StatusKawin - Status perkawinan
var StatusKawin = map[int]string{
	1: "BELUM KAWIN",
	2: "KAWIN",
	3: "CERAI HIDUP",
	4: "CERAI MATI",
}

// HubunganKeluarga - Status hubungan dalam keluarga
var HubunganKeluarga = map[int]string{
	1:  "KEPALA KELUARGA",
	2:  "SUAMI",
	3:  "ISTRI",
	4:  "ANAK",
	5:  "MENANTU",
	6:  "CUCU",
	7:  "ORANGTUA",
	8:  "MERTUA",
	9:  "FAMILI LAIN",
	10: "PEMBANTU",
	11: "LAINNYA",
}

// Pendidikan - Tingkat pendidikan
var Pendidikan = map[int]string{
	1:  "TIDAK / BELUM SEKOLAH",
	2:  "BELUM TAMAT SD/SEDERAJAT",
	3:  "TAMAT SD / SEDERAJAT",
	4:  "SLTP / SEDERAJAT",
	5:  "SLTA / SEDERAJAT",
	6:  "DIPLOMA I/II",
	7:  "AKADEMI/DIPLOMA III/S. MUDA",
	8:  "DIPLOMA IV/STRATA I",
	9:  "STRATA II",
	10: "STRATA III",
}

// WargaNegara - Kewarganegaraan
var WargaNegara = map[int]string{
	1: "WNI",
	2: "WNA",
}

// Pekerjaan - Daftar pekerjaan
var Pekerjaan = map[int]string{
	1:  "BELUM/TIDAK BEKERJA",
	2:  "MENGURUS RUMAH TANGGA",
	3:  "PELAJAR/MAHASISWA",
	4:  "PENSIUNAN",
	5:  "PEGAWAI NEGERI SIPIL (PNS)",
	6:  "TENTARA NASIONAL INDONESIA (TNI)",
	7:  "KEPOLISIAN RI (POLRI)",
	8:  "PERDAGANGAN",
	9:  "PETANI/PEKEBUN",
	10: "PETERNAK",
	11: "NELAYAN/PERIKANAN",
	12: "INDUSTRI",
	13: "KONSTRUKSI",
	14: "TRANSPORTASI",
	15: "KARYAWAN SWASTA",
	16: "KARYAWAN BUMN",
	17: "KARYAWAN BUMD",
	18: "KARYAWAN HONORER",
	19: "BURUH HARIAN LEPAS",
	20: "BURUH TANI/PERKEBUNAN",
	21: "BURUH NELAYAN/PERIKANAN",
	22: "BURUH PETERNAKAN",
	23: "PEMBANTU RUMAH TANGGA",
	24: "TUKANG CUKUR",
	25: "TUKANG LISTRIK",
	26: "TUKANG BATU",
	27: "TUKANG KAYU",
	28: "TUKANG SOL SEPATU",
	29: "TUKANG LAS/PANDAI BESI",
	30: "TUKANG JAHIT",
	31: "TUKANG GIGI",
	32: "PENATA RIAS",
	33: "PENATA BUSANA",
	34: "PENATA RAMBUT",
	35: "MEKANIK",
	36: "SENIMAN",
	37: "TABIB",
	38: "PARAJI",
	39: "PERANCANG BUSANA",
	40: "PENTERJEMAH",
	41: "IMAM MASJID",
	42: "PENDETA",
	43: "PASTOR",
	44: "WARTAWAN",
	45: "USTADZ/MUBALIGH",
	46: "JURU MASAK",
	47: "PROMOTOR ACARA",
	48: "ANGGOTA DPR-RI",
	49: "ANGGOTA DPD",
	50: "ANGGOTA BPK",
	51: "PRESIDEN",
	52: "WAKIL PRESIDEN",
	53: "ANGGOTA MAHKAMAH KONSTITUSI",
	54: "ANGGOTA KABINET KEMENTERIAN",
	55: "DUTA BESAR",
	56: "GUBERNUR",
	57: "WAKIL GUBERNUR",
	58: "BUPATI",
	59: "WAKIL BUPATI",
	60: "WALIKOTA",
	61: "WAKIL WALIKOTA",
	62: "ANGGOTA DPRD PROVINSI",
	63: "ANGGOTA DPRD KABUPATEN/KOTA",
	64: "DOSEN",
	65: "GURU",
	66: "PILOT",
	67: "PENGACARA",
	68: "NOTARIS",
	69: "ARSITEK",
	70: "AKUNTAN",
	71: "KONSULTAN",
	72: "DOKTER",
	73: "BIDAN",
	74: "PERAWAT",
	75: "APOTEKER",
	76: "PSIKIATER/PSIKOLOG",
	77: "PENYIAR TELEVISI",
	78: "PENYIAR RADIO",
	79: "PELAUT",
	80: "PENELITI",
	81: "SOPIR",
	82: "PIALANG",
	83: "PARANORMAL",
	84: "PEDAGANG",
	85: "PERANGKAT DESA",
	86: "KEPALA DESA",
	87: "BIARAWATI",
	88: "WIRASWASTA",
	89: "LAINNYA",
}

// Status - Status aktif/tidak aktif
var Status = map[int]string{
	1: "AKTIF",
	2: "TIDAK AKTIF",
}

// Domisili - Status domisili
var Domisili = map[int]string{
	1: "PENDUDUK TETAP",
	2: "PENDUDUK DOMISILI",
}

// Helper functions untuk mendapatkan nilai string
func GetKelamin(code int) string {
	if val, ok := JenisKelamin[code]; ok {
		return val
	}
	return ""
}

func GetJK(code int) string {
	if val, ok := JK[code]; ok {
		return val
	}
	return ""
}

func GetAgama(code int) string {
	if val, ok := Agama[code]; ok {
		return val
	}
	return ""
}

func GetStatusKawin(code int) string {
	if val, ok := StatusKawin[code]; ok {
		return val
	}
	return ""
}

func GetHubunganKeluarga(code int) string {
	if val, ok := HubunganKeluarga[code]; ok {
		return val
	}
	return ""
}

func GetPendidikan(code int) string {
	if val, ok := Pendidikan[code]; ok {
		return val
	}
	return ""
}

func GetWargaNegara(code int) string {
	if val, ok := WargaNegara[code]; ok {
		return val
	}
	return ""
}

func GetPekerjaan(code int) string {
	if val, ok := Pekerjaan[code]; ok {
		return val
	}
	return ""
}

func GetStatus(code int) string {
	if val, ok := Status[code]; ok {
		return val
	}
	return ""
}

func GetDomisili(code int) string {
	if val, ok := Domisili[code]; ok {
		return val
	}
	return ""
}
