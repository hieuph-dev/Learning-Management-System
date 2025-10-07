package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateSlug tạo slug từ title
func GenerateSlug(title string) string {
	// Chuyển về lowercase
	slug := strings.ToLower(title)

	// Loại bỏ dấu tiếng Việt
	slug = removeVietnameseTones(slug)

	// Thay thế khoảng trắng và ký tự đặc biệt bằng dấu gạch ngang
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Loại bỏ dấu gạch ngang ở đầu và cuối
	slug = strings.Trim(slug, "-")

	// Loại bỏ các dấu gạch ngang liên tiếp
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}

// removeVietnameseTones loại bỏ dấu tiếng Việt
func removeVietnameseTones(s string) string {
	// Map các ký tự có dấu sang không dấu
	replacements := map[rune]string{
		'à': "a", 'á': "a", 'ạ': "a", 'ả': "a", 'ã': "a",
		'â': "a", 'ầ': "a", 'ấ': "a", 'ậ': "a", 'ẩ': "a", 'ẫ': "a",
		'ă': "a", 'ằ': "a", 'ắ': "a", 'ặ': "a", 'ẳ': "a", 'ẵ': "a",
		'è': "e", 'é': "e", 'ẹ': "e", 'ẻ': "e", 'ẽ': "e",
		'ê': "e", 'ề': "e", 'ế': "e", 'ệ': "e", 'ể': "e", 'ễ': "e",
		'ì': "i", 'í': "i", 'ị': "i", 'ỉ': "i", 'ĩ': "i",
		'ò': "o", 'ó': "o", 'ọ': "o", 'ỏ': "o", 'õ': "o",
		'ô': "o", 'ồ': "o", 'ố': "o", 'ộ': "o", 'ổ': "o", 'ỗ': "o",
		'ơ': "o", 'ờ': "o", 'ớ': "o", 'ợ': "o", 'ở': "o", 'ỡ': "o",
		'ù': "u", 'ú': "u", 'ụ': "u", 'ủ': "u", 'ũ': "u",
		'ư': "u", 'ừ': "u", 'ứ': "u", 'ự': "u", 'ử': "u", 'ữ': "u",
		'ỳ': "y", 'ý': "y", 'ỵ': "y", 'ỷ': "y", 'ỹ': "y",
		'đ': "d",
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			result.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}

	// Normalize unicode
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)

	normalized, _, _ := transform.String(t, result.String())
	return normalized
}

// GenerateUniqueSlug tạo slug duy nhất (thêm số nếu trùng)
func GenerateUniqueSlug(baseSlug string, existingSlugCheck func(string) bool) string {
	slug := baseSlug
	counter := 1

	for existingSlugCheck(slug) {
		slug = baseSlug + "-" + string(rune(counter))
		counter++
	}

	return slug
}
