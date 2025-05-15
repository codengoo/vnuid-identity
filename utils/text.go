package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func RemoveVietnameseTones(input string) string {
	// Chuẩn hóa chuỗi về dạng NFD (phân tách chữ và dấu)
	t := norm.NFD.String(input)

	var b strings.Builder
	for _, r := range t {
		// Loại bỏ các ký tự dấu (thuộc class Mn - Mark, nonspacing)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// Chuyển đ → d, Đ → D
		if r == 'đ' {
			r = 'd'
		} else if r == 'Đ' {
			r = 'D'
		}
		b.WriteRune(r)
	}

	// Loại bỏ ký tự đặc biệt và chuẩn hóa khoảng trắng
	result := b.String()
	result = strings.TrimSpace(result)
	result = strings.Join(strings.Fields(result), " ") // chuẩn hóa khoảng trắng

	return result
}
