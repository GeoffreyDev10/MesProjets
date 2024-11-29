package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	params := os.Args
	long := 0

	for index := range params {
		long = index
	}

	for i := long; i >= 1; i-- {
		for _, j := range params[i] {
			z01.PrintRune(j)
		}
		z01.PrintRune('\n')
	}
}
