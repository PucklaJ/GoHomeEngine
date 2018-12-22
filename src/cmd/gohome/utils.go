package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), Env...)
	return cmd.Run()
}

func ConsoleRead() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func WorkingDir() string {
	str, _ := os.Getwd()
	return str + "/"
}

func FileExists(fn string) bool {
	_, err := os.Stat(fn)
	return !os.IsNotExist(err)
}

func PackageName() string {
	wd := WorkingDir()
	return wd[strings.LastIndex(wd[:len(wd)-1], "/")+1 : len(wd)-1]
}
