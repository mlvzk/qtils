package util

func LeftPad(s string, padStr string, pLen int) string {
	for i := len(s); i < pLen; i++ {
		s = padStr + s
	}
	return s
}
func RightPad(s string, padStr string, pLen int) string {
	for i := len(s); i < pLen; i++ {
		s = s + padStr
	}
	return s
}
