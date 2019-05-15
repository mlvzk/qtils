package util

// TODO: remove dependency on `strings`
import "strings"

func LeftPad(s string, padStr string, pLen int) string {
	if len(s) >= pLen {
		return s
	}

	return strings.Repeat(padStr, pLen-len(s)) + s
}
func RightPad(s string, padStr string, pLen int) string {
	if len(s) >= pLen {
		return s
	}

	return s + strings.Repeat(padStr, pLen-len(s))
}
