package utils

import (
	"strings"
	"testing"
)

func TestRandomHex(t *testing.T) {
	s, err := randomHex(16)
	if err != nil {
		t.Fatalf("randomHex error: %v", err)
	}
	if len(s) != 32 { // 16 byte -> 32 karakter hex
		t.Errorf("panjang = %d, want 32", len(s))
	}

	s2, _ := randomHex(16)
	if s == s2 {
		t.Error("dua panggilan randomHex menghasilkan nilai sama")
	}
}

func TestDetectContentType(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if got := detectContentType(png); got != "image/png" {
		t.Errorf("PNG terdeteksi sebagai %q", got)
	}

	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	if got := detectContentType(jpeg); got != "image/jpeg" {
		t.Errorf("JPEG terdeteksi sebagai %q", got)
	}

	// Konten teks tidak boleh dianggap gambar.
	if got := detectContentType([]byte("halo dunia bukan gambar")); strings.HasPrefix(got, "image/") {
		t.Errorf("teks salah terdeteksi sebagai gambar: %q", got)
	}
}
