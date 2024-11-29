package piscine

func IsUpper(str string) bool {
	for ch := 0; ch < len(str); ch++ {
		if str[ch] < 'A' || str[ch] > 'Z' {
			return false
		}
	}

	return true
}
