package utils

func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	runes := []rune(s)

	if len(runes) <= maxLen {
		return s
	}

	return string(runes[:maxLen]) + "..."
}
