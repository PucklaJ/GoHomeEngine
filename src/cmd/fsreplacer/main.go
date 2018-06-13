package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 3 || len(argsWithoutProg) > 3 {
		fmt.Println("Only 3 arguments are supported: file string_to_replace string_which_gets_into_file")
		os.Exit(1)
	}

	if len(argsWithoutProg[1]) == 0 || len(argsWithoutProg[2]) == 0 {
		fmt.Println("The length of the relpacement strings have to be at least one")
		os.Exit(1)
	}

	contentsBytes, err := ioutil.ReadFile(argsWithoutProg[0])
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	contents := string(contentsBytes)

	newContents := strings.Replace(contents, argsWithoutProg[1], argsWithoutProg[2], -1)

	err1 := ioutil.WriteFile(argsWithoutProg[0], []byte(newContents), os.ModeType)
	if err1 != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	os.Exit(0)
}
