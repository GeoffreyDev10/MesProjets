package piscine

func StrRev(s string) string {
	reverse := ""
	for _, c := range s {
		reverse = string(c) + reverse
	}
	return reverse
}
