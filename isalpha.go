package piscine

func IsAlpha(s string) bool {
	for ch := 0; ch < len(s); ch++ {
		if s[ch] < 'a' || s[ch] > 'z' {
			if s[ch] < 'A' || s[ch] > 'Z' {
				if s[ch] < '0' || s[ch] > '9' {
					return false
				}
			}
		}
	}
	return true
}
