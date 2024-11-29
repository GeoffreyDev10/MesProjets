package piscine

func NRune(s string, n int) rune {
	ptarg := []rune(s)
	len := 0
	for index := range ptarg {
		len = index
	}
	if n-1 >= 0 && n-1 <= len {
		return ptarg[n-1]
	}
	return 0
}
