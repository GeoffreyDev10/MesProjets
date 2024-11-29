package piscine

func IterativePower(nb int, power int) int {
	if power < 0 {
		return 0
	} else {
		if power == 0 {
			return 1
		}
	}
	a := nb
	for i := power; i > 1; i-- {
		a = nb * a
	}
	return a
}
