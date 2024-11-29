package piscine

func ToUpper(s string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		arg := s[i]
		if arg >= 'a' && arg <= 'z' {
			arg = arg - 32
		}
		result += string(arg)
	}
	return result
}
