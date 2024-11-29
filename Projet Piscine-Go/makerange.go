package piscine

func MakeRange(min, max int) []int {
	if max <= min {
		return nil
	}

	var result []int = make([]int, max-min)
	for i, tempA := min, 0; i < max; i++ {
		result[tempA] = i
		tempA++
	}
	return result
}
