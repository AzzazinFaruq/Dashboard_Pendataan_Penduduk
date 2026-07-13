package config

import "testing"

func TestGetPekerjaan(t *testing.T) {
	cases := map[int]string{
		1:   "BELUM/TIDAK BEKERJA",
		10:  "PETERNAK",
		88:  "WIRASWASTA", // sebelumnya hilang (map hanya sampai 10), regression guard
		89:  "LAINNYA",
		0:   "", // kode tidak dikenal
		999: "",
	}
	for code, want := range cases {
		if got := GetPekerjaan(code); got != want {
			t.Errorf("GetPekerjaan(%d) = %q, want %q", code, got, want)
		}
	}
}

func TestGetKelamin(t *testing.T) {
	if got := GetKelamin(1); got != "LAKI-LAKI" {
		t.Errorf("GetKelamin(1) = %q, want LAKI-LAKI", got)
	}
	if got := GetKelamin(2); got != "PEREMPUAN" {
		t.Errorf("GetKelamin(2) = %q, want PEREMPUAN", got)
	}
	if got := GetKelamin(42); got != "" {
		t.Errorf("GetKelamin(42) = %q, want empty", got)
	}
}

func TestGetStatusDanDomisili(t *testing.T) {
	if GetStatus(1) != "AKTIF" || GetStatus(2) != "TIDAK AKTIF" {
		t.Error("GetStatus mengembalikan nilai tak terduga")
	}
	if GetDomisili(1) != "PENDUDUK TETAP" {
		t.Error("GetDomisili(1) salah")
	}
}
