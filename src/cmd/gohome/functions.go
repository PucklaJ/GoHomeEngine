package main

import (
	"fmt"
	"os"
)

const CONFIG_FILE_NAME = ".gohome.config"

func HandleConfigFile() {
	CustomValues = make(map[string]string)
	old_CustomValues = make(map[string]string)
	wd := WorkingDir()
	fn := wd + CONFIG_FILE_NAME
	if FileExists(fn) {
		file, err := os.Open(fn)
		if err != nil {
			fmt.Println("Failed to open", CONFIG_FILE_NAME, ":", err)
			os.Exit(1)
		}
		readVariables(file)
		file.Close()
	} else {
		file, err := os.Create(fn)
		if err != nil {
			fmt.Println("Failed to create", CONFIG_FILE_NAME, ":", err)
			os.Exit(1)
		}
		writeVariables(file)
		file.Close()
	}
}

func HandleArguments() {
	if len(os.Args) == 1 {
		fmt.Println("No Arguments")
		return
	}

	for _, arg := range os.Args[1:] {
		if isCommandArg(arg) {
			COMMAND = arg
			continue
		}
		if isValueArg(arg) {
			if processValueArg(arg) {
				continue
			}
		}
		if isConfigArg(arg) {
			VAR_CONFIG = arg
			continue
		}
		if arg == "-a" || arg == "--all" {
			continue
		}

		fmt.Println("Invalid argument:", arg)
		os.Exit(1)
	}

	if COMMAND == "" {
		fmt.Println("No command specified: build|install|run|generate|clean|env|set|reset")
		WriteConfigFile()
		os.Exit(1)
	}
}

func ExecuteCommands() {
	build := getBuild()
	if build == nil {
		fmt.Println("OS not supported:", VAR_OS)
		WriteConfigFile()
		os.Exit(1)
	}

	var success = true

	switch COMMAND {
	case "build":
		if !build.IsGenerated() || valuesChanged() {
			build.Generate()
		}
		success = build.Build()
	case "install":
		if !build.IsGenerated() || valuesChanged() {
			build.Generate()
		}
		success = build.Install()
	case "run":
		if !build.IsGenerated() || valuesChanged() {
			build.Generate()
		}
		success = build.Build()
		if success {
			success = build.Run()
		}
	case "generate":
		build.Generate()
	case "clean":
		build.Clean()
	case "env":
		build.Env()
	case "reset":
		resetParameters()
	case "export":
		if !build.IsGenerated() || valuesChanged() {
			build.Generate()
		}
		build.Build()
		build.Export()
	}

	if !success {
		WriteConfigFile()
		os.Exit(1)
	}
}

func WriteConfigFile() {
	file, err := os.Create(WorkingDir() + CONFIG_FILE_NAME)
	if err != nil {
		fmt.Println("Failed to open", CONFIG_FILE_NAME, ":", err)
		os.Exit(1)
	}
	writeVariables(file)
	file.Close()
}
