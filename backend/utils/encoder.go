package utils

import (
	"sayhi/backend/models"
	"unicode/utf8"
)

// CountChars 根据编码类型计算字符数
func CountChars(text string, encoding models.EncodingType) int {
	switch encoding {
	case models.EncodingASCII:
		return countASCII(text)
	case models.EncodingZawgyi:
		return countZawgyi(text)
	case models.EncodingUnicode:
		return utf8.RuneCountInString(text)
	case models.EncodingOther:
		return len([]byte(text))
	default:
		return utf8.RuneCountInString(text)
	}
}

// countASCII 计算ASCII字符数（只计算ASCII字符）
func countASCII(text string) int {
	count := 0
	for _, r := range text {
		if r < 128 {
			count++
		}
	}
	return count
}

// countZawgyi 计算Zawgyi字符数（按字节计算，Zawgyi通常需要更多字节）
func countZawgyi(text string) int {
	// Zawgyi编码通常每个字符占用2-3字节
	// 这里简化为按字节数计算
	return len([]byte(text))
}

// IsExceeded 检查是否超出限制
func IsExceeded(charCount int, limit int) bool {
	return charCount > limit
}
