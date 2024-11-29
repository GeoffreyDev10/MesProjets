package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	for _, arg := range os.Args[1:] {
		n := 0
		for _, char := range arg {
			if char < '0' || char > '9' {
				n = -1
				break
			}
			n = n*10 + int(char-'0')
		}
		if n >= 1 && n <= 26 {
		}
	}
	z01.PrintRune('\n')
}
