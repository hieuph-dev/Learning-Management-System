package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

// GenerateSecureToken tạo token ngẫu nhiên an toàn (cryptographically secure).
// Trả về chuỗi base64 URL-safe (không chứa +, /, =).
func GenerateSecureToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid token length: %d", length)
	}

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken băm token bằng SHA-256 để lưu vào DB.
// Best practice: không lưu token raw trong DB, chỉ lưu hash.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// GenerateResetCode tạo mã số 6 chữ số an toàn (crypto-secure).
func GenerateResetCode() (string, error) {
	max := big.NewInt(1000000) // [0, 999999]
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// GeneratePasswordResetToken tạo token kép (URL token + mã nhập tay).
// secureToken: dùng trong link gửi qua email
// readableCode: mã ngắn user nhập thủ công
// hashToken: nên lưu hashToken vào DB để so sánh
func GeneratePasswordResetToken() (secureToken, readableCode, hashToken string, err error) {
	secureToken, err = GenerateSecureToken(32) // 32 bytes → ~43 ký tự base64
	if err != nil {
		return "", "", "", err
	}

	readableCode, err = GenerateResetCode()
	if err != nil {
		return "", "", "", err
	}

	hashToken = HashToken(secureToken)
	return secureToken, readableCode, hashToken, nil
}

// IsTokenExpired kiểm tra token có hết hạn chưa.
func IsTokenExpired(expiresAt time.Time) bool {
	return !time.Now().Before(expiresAt) // true nếu now >= expiresAt
}

// GetResetTokenExpiry trả về thời điểm hết hạn (mặc định 1 giờ).
func GetResetTokenExpiry() time.Time {
	return time.Now().UTC().Add(1 * time.Hour)
}
