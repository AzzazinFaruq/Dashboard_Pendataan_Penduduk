package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// allowedImageTypes adalah MIME type gambar yang diizinkan untuk diupload.
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

// maxImageSize adalah batas ukuran file gambar (5 MB).
const maxImageSize = 5 * 1024 * 1024

// ValidateAndBuildImagePath memvalidasi file gambar yang diupload berdasarkan
// content-type asli (bukan ekstensi nama file) dan ukurannya, lalu mengembalikan
// nama file acak yang aman untuk disimpan pada direktori dir.
//
// Mengembalikan error jika tipe tidak didukung atau ukuran melebihi batas.
func ValidateAndBuildImagePath(dir string, file *multipart.FileHeader) (string, error) {
	if file.Size > maxImageSize {
		return "", fmt.Errorf("ukuran file terlalu besar (maks 5MB)")
	}

	// Deteksi tipe dari isi file, bukan dari nama file yang bisa dipalsukan.
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membaca file")
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	contentType := detectContentType(buf[:n])

	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return "", fmt.Errorf("format file tidak didukung (hanya jpg/png)")
	}

	// Nama file acak untuk mencegah path traversal & tabrakan nama.
	name, err := randomHex(16)
	if err != nil {
		return "", fmt.Errorf("gagal membuat nama file")
	}

	// Bersihkan dir dan gabungkan dengan nama file aman.
	cleanDir := filepath.Clean(dir)
	return filepath.ToSlash(filepath.Join(cleanDir, name+ext)), nil
}

// detectContentType mendeteksi MIME type dari byte awal file.
func detectContentType(data []byte) string {
	ct := http.DetectContentType(data)
	// http.DetectContentType kadang menambahkan "; charset=..."; ambil bagian tipe saja.
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}
	return strings.TrimSpace(ct)
}

// randomHex menghasilkan string heksadesimal acak sepanjang n byte.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
