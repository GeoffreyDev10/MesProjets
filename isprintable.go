package piscine

func IsPrintable(str string) bool {
	ch := str
	ln := 0
	for i := range ch {
		ln = i
	}

	for i := 0; i <= ln; i++ {
		if ch[i] < 32 {
			return false
		}
	}
	return true
}
