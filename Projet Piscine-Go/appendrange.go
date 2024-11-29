package piscine

func AppendRange(min, max int) []int {
	if min >= max {
		var result []int = nil
		return result
	} else {
		result := []int{}
		for i := min; i < max; i++ {
			result = append(result, i)
		}
		return result
	}
}
