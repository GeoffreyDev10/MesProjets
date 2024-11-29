package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	args := os.Args[1:]
	for _, res := range args {
		for _, char := range res {
			z01.PrintRune(char)
		}
		z01.PrintRune('\n')
	}
}
