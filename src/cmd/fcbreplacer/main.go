package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("Error: No Arguments")
		os.Exit(1)
	}

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 1 {
		fmt.Println("Error: Only one argument is supported")
		os.Exit(1)
	}

	str := argsWithoutProg[0]

	if str[0] >= 97 && str[0] <= 122 {
		str = strings.Replace(str, string(str[0]), string(str[0]-97+65), 1)
	}

	fmt.Print(str)

	os.Exit(0)
}
