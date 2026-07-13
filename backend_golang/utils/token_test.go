package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	// GenerateJWT/ValidateJWT membaca JWT_SECRET dari environment.
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Exit(m.Run())
}

func TestGenerateAndValidateJWT(t *testing.T) {
	token, err := GenerateJWT(42, time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	parsed, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("token seharusnya valid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims bukan MapClaims")
	}
	if sub, _ := claims["sub"].(float64); int(sub) != 42 {
		t.Errorf("sub = %v, want 42", claims["sub"])
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	token, err := GenerateJWT(1, time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	// Ganti secret setelah token dibuat -> validasi harus gagal.
	os.Setenv("JWT_SECRET", "secret-yang-berbeda")
	defer os.Setenv("JWT_SECRET", "test-secret-key")

	if _, err := ValidateJWT(token); err == nil {
		t.Error("ValidateJWT seharusnya gagal untuk secret yang berbeda")
	}
}

func TestGenerateJWTDurationApplied(t *testing.T) {
	token, err := GenerateJWT(1, time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	parsed, _ := ValidateJWT(token)
	claims := parsed.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))

	// exp harus berada di sekitar 1 jam dari sekarang (toleransi 1 menit).
	expected := time.Now().Add(time.Hour).Unix()
	if diff := exp - expected; diff < -60 || diff > 60 {
		t.Errorf("exp menyimpang %d detik dari ekspektasi 1 jam", diff)
	}
}
