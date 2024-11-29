package piscine

func LastRune(s string) rune {
	ptarg := []rune(s)
	return ptarg[len(s)-1]
}
