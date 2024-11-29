package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	programPath := os.Args[0]

	programName := programPath
	for i := len(programPath) - 1; i >= 0; i-- {
		if programPath[i] == '/' {
			programName = programPath[i+1:]
			break
		}
	}
	for _, cname := range programName {
		z01.PrintRune(cname)
	}
	z01.PrintRune('\n')
}
