package piscine

func IsNumeric(s string) bool {
	for ch := 0; ch < len(s); ch++ {
		if s[ch] < '0' || s[ch] > '9' {
			return false
		}
	}
	return true
}
