package main

import (
	"fmt"
	"os"
	"strings"
)

func getNumberAsWord(number byte) string {
	number -= 48
	switch number {
	case 0:
		return "Zero"
	case 1:
		return "One"
	case 2:
		return "Two"
	case 3:
		return "Three"
	case 4:
		return "Four"
	case 5:
		return "Five"
	case 6:
		return "Six"
	case 7:
		return "Seven"
	case 8:
		return "Eight"
	case 9:
		return "Nine"
	default:
		return "Nan"
	}
}

func main() {
	if len(os.Args) == 0 {
		fmt.Println("Error: No Arguments")
		os.Exit(1)
	}
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		fmt.Println("Error: Only one argument is allowed")
		os.Exit(1)
	}

	str := argsWithoutProg[0]

	if str[0] >= 48 && str[0] <= 57 {
		newStr := strings.Replace(str, string(str[0]), getNumberAsWord(str[0]), 1)
		fmt.Print(newStr)
	} else {
		fmt.Print(str)
	}

	os.Exit(0)
}
