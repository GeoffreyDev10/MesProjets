package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arg1 := os.Args[1:]
	for i := 1; i < len(arg1); i++ {
		key := arg1[i]
		j := i - 1

		for j >= 0 && arg1[j] > key {
			arg1[j+1] = arg1[j]
			j--
		}
		arg1[j+1] = key
	}

	for _, res := range arg1 {
		for _, char := range res {
			z01.PrintRune(char)
		}
		z01.PrintRune('\n')
	}
}
