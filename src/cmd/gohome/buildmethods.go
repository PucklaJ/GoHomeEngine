package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func (*DesktopBuild) build(str string) bool {
	var varos, vararch string
	if VAR_OS == "runtime" {
		varos = runtime.GOOS
	} else {
		varos = VAR_OS
	}
	if VAR_ARCH == "runtime" {
		vararch = runtime.GOARCH
	} else {
		vararch = VAR_ARCH
	}

	Env = []string{
		"GOOS=" + varos,
		"GOARCH=" + vararch,
	}

	var err error
	if VAR_CONFIG == "DEBUG" {
		Env = append(Env, []string{
			"CGO_FLAGS=-g",
			"CGO_LDFLAGS=-g",
			"CGO_CXXFLAGS=-g",
		}...)
		err = ExecCommand("go", str, "-v")
	} else {
		Env = append(Env, []string{
			"CGO_FLAGS=-O3",
			"CGO_LDFLAGS=-O3",
			"CGO_CXXFLAGS=-O3",
		}...)
		err = ExecCommand("go", str, "-v")
		if err == nil {
			ExecCommand("strip", "-s", "-x", "--strip-unneeded", PackageName())
		}
	}

	return err == nil
}

func (this *DesktopBuild) Build() bool {
	return this.build("build")
}
func (this *DesktopBuild) Install() bool {
	return this.build("install")
}

func (this *DesktopBuild) generateMain() (str string) {
	str += "package main\n\n"
	str += "import (\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/" + VAR_FRAME + "\"\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/gohome\"\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/" + VAR_RENDER + "\"\n"
	str += ")\n\n"
	str += "func main() {\n"
	var frame string
	if VAR_FRAME == "GTK" {
		frame = "&framework." + VAR_FRAME + "Framework{\n\t\tUseWholeWindowAsGLArea: " + strconv.FormatBool(this.gtkwholewindow) + ",\n\t\tMenuBarFix: " + strconv.FormatBool(this.gtkmenubar) + ",\n\t}"
	} else {
		frame = "&framework." + VAR_FRAME + "Framework{}"
	}

	str += "\tgohome.MainLop.Run(" + frame + ",&renderer." + VAR_RENDER + "Renderer{}," + strconv.FormatInt(int64(this.width), 10) + "," + strconv.FormatInt(int64(this.height), 10) + ",\"" + this.title + "\",&" + VAR_START + "{})\n"
	str += "}\n"

	return
}

func (this *DesktopBuild) Generate() {
	if VAR_FRAME == "GTK" && VAR_RENDER != "OpenGL" {
		fmt.Println(VAR_FRAME, "is not compatible with", VAR_RENDER)
		os.Exit(1)
	}

	this.title = CustomValues["TITLE"]
	if str, ok := CustomValues["WIDTH"]; ok {
		i, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			this.width = int(i)
		}
	}
	if str, ok := CustomValues["HEIGHT"]; ok {
		i, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			this.height = int(i)
		}
	}

	if this.title == "" {
		fmt.Print("Title: ")
		this.title = ConsoleRead()
		CustomValues["TITLE"] = this.title
	}
	if this.width == 0 {
		fmt.Print("Width: ")
		i, err := strconv.ParseInt(ConsoleRead(), 10, 32)
		if err != nil {
			this.width = 1280
		} else {
			this.width = int(i)
		}
		CustomValues["WIDTH"] = strconv.FormatInt(int64(this.width), 10)
	}
	if this.height == 0 {
		fmt.Print("Height: ")
		i, err := strconv.ParseInt(ConsoleRead(), 10, 32)
		if err != nil {
			this.height = 720
		} else {
			this.height = int(i)
		}
		CustomValues["HEIGHT"] = strconv.FormatInt(int64(this.height), 10)
	}

	if VAR_FRAME == "GTK" {
		var err error
		var str string
		var ok bool
		if str, ok = CustomValues["USEWHOLEWINDOWASGLAREA"]; ok {
			this.gtkwholewindow, err = strconv.ParseBool(str)
		}

		if !ok || err != nil {
			fmt.Print("UseWholeWindowAsGLArea: ")
			this.gtkwholewindow, err = strconv.ParseBool(ConsoleRead())
			if err != nil {
				this.gtkwholewindow = true
			}
		}
		err = nil
		if str, ok = CustomValues["MENUBARFIX"]; ok {
			this.gtkmenubar, err = strconv.ParseBool(str)
		}

		if !ok || err != nil {
			fmt.Print("MenuBarFix: ")
			this.gtkmenubar, err = strconv.ParseBool(ConsoleRead())
			if err != nil {
				this.gtkmenubar = false
			}
		}

		CustomValues["USEWHOLEWINDOWASGLAREA"] = strconv.FormatBool(this.gtkwholewindow)
		CustomValues["MENUBARFIX"] = strconv.FormatBool(this.gtkmenubar)
	}

	if VAR_START == "" {
		fmt.Print("StartScene: ")
		VAR_START = ConsoleRead()
	}
	str := this.generateMain()
	file, err := os.Create(WorkingDir() + "main.go")
	if err != nil {
		fmt.Println("Failed to generate main.go:", err)
		os.Exit(1)
	}
	file.WriteString(str)
	file.Close()
}
func (*DesktopBuild) IsGenerated() bool {
	return FileExists(WorkingDir() + "main.go")
}
func (*DesktopBuild) Run() bool {
	pack := PackageName()
	if runtime.GOOS == "windows" {
		pack += ".exe"
	}
	err := ExecCommand("./" + pack)
	return err == nil
}
func (*DesktopBuild) Export() {

}
func (*DesktopBuild) Clean() {
	ExecCommand("go", "clean", "-r", "--cache")
	ExecCommand("rm", "-f", "main.go")
}

func (*DesktopBuild) Env() {
	fmt.Println("OS=" + VAR_OS)
	fmt.Println("ARCH=" + VAR_ARCH)
	fmt.Println("FRAME=" + VAR_FRAME)
	fmt.Println("RENDER=" + VAR_RENDER)
	fmt.Println("START=" + VAR_START)
	fmt.Println("CONFIG=" + VAR_CONFIG)
	for k, v := range CustomValues {
		fmt.Println(k + "=" + v)
	}
	var all = false
	for _, arg := range os.Args {
		if arg == "-a" || arg == "--all" {
			all = true
		}
	}

	if all {
		ExecCommand("go", "env")
	}
}

func (*AndroidBuild) Build() bool {
	return true
}
func (*AndroidBuild) Install() bool {
	return true
}
func (*AndroidBuild) Generate() {

}
func (*AndroidBuild) IsGenerated() bool {
	return false
}
func (*AndroidBuild) Run() bool {
	return true
}
func (*AndroidBuild) Export() {

}
func (*AndroidBuild) Clean() {

}
func (*AndroidBuild) Env() {

}
