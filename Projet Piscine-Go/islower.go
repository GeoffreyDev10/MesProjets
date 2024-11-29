package piscine

func IsLower(str string) bool {
	for ch := 0; ch < len(str); ch++ {
		if str[ch] < 'a' || str[ch] > 'z' {
			return false
		}
	}

	return true
}
