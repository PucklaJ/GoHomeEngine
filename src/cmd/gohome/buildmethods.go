package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
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

	var ldflags string
	if VAR_CONFIG == "RELEASE" {
		ldflags = "-ldflags=-s -w"
		if runtime.GOOS == "windows" {
			ldflags += " -extldflags=-Wl,--subsystem,windows"
		}
	}
	if COMMAND == "export" {
		if VAR_CONFIG == "DEBUG" {
			ldflags = "-ldflags="
		} else {
			ldflags += " "
		}

		if runtime.GOOS == "linux" {
			ldflags += "-extldflags=-Wl,-z,origin,-rpath=$ORIGIN/lib"
		}
	}

	var err error
	if VAR_CONFIG == "DEBUG" {
		Env = append(Env, []string{
			"CGO_FLAGS=-g",
			"CGO_LDFLAGS=-g",
			"CGO_CXXFLAGS=-g",
		}...)
		err = ExecCommand("go", str, "-v", ldflags)
	} else {
		Env = append(Env, []string{
			"CGO_FLAGS=-O3",
			"CGO_LDFLAGS=-O3",
			"CGO_CXXFLAGS=-O3",
		}...)
		err = ExecCommand("go", str, "-v", ldflags)
	}

	return err == nil
}

func (this *DesktopBuild) Build() bool {
	return this.build("build")
}
func (this *DesktopBuild) Install() bool {
	return this.build("install")
}

func generateMain(forandroid bool) (str string) {
	str += "package main\n\n"
	if forandroid {
		str += "import \"C\"\n\n"
	}
	str += "import (\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/" + VAR_FRAME + "\"\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/gohome\"\n"
	str += "\t\"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/" + VAR_RENDER + "\"\n"
	str += ")\n\n"
	str += "func main() {\n"
	var frame string
	if VAR_FRAME == "GTK" {
		frame = "&framework." + VAR_FRAME + "Framework{\n\t\tUseWholeWindowAsGLArea: " + CustomValues["USEWHOLEWINDOWASGLAREA"] + ",\n\t\tMenuBarFix: " + CustomValues["MENUBARFIX"] + ",\n\t}"
	} else {
		frame = "&framework." + VAR_FRAME + "Framework{}"
	}

	str += "\tgohome.MainLop.Run(" + frame + ",&renderer." + VAR_RENDER + "Renderer{}," + CustomValues["WIDTH"] + "," + CustomValues["HEIGHT"] + ",\"" + CustomValues["TITLE"] + "\",&" + VAR_START + "{})\n"
	str += "}\n"

	if forandroid {
		str += "\n//export SDL_main\n"
		str += "func SDL_main() {\n"
		str += "\tmain()\n"
		str += "}\n"
	}

	return
}

func (this *DesktopBuild) Generate() {
	if VAR_FRAME == "JS" {
		fmt.Println("Desktop is not compatible with JS")
		AssertValue(&VAR_FRAME, "JS", "Framework")
	}
	if VAR_FRAME == "GTK" && VAR_RENDER != "OpenGL" {
		fmt.Println(VAR_FRAME, "is not compatible with", VAR_RENDER)
		os.Exit(1)
	}
	if VAR_RENDER == "WebGL" {
		fmt.Println("Desktop is not compatible with WebGL")
		fmt.Println("(1) OpenGL")
		fmt.Println("(2) OpenGLES2")
		fmt.Print("Which renderer: ")
		render := ConsoleReadi()
		switch render {
		case 1:
			VAR_RENDER = "OpenGL"
		case 2:
			VAR_RENDER = "OpenGLES2"
		}
	}

	if VAR_RENDER == "OpenGLES3" || VAR_RENDER == "OpenGLES31" {
		fmt.Println("Desktop is only compatible with OpenGL or OpenGLES2!")
		fmt.Print("(1) OpenGL\n(2)OpenGLES2\nChoose one: ")
		render := ConsoleReadi()
		switch render {
		case 2:
			VAR_RENDER = "OpenGLES2"
		default:
			VAR_RENDER = "OpenGL"
		}
	}

	GetCustomValue("TITLE")
	GetCustomValuei("WIDTH", 1280)
	GetCustomValuei("HEIGHT", 720)

	if VAR_FRAME == "GTK" {
		this.gtkwholewindow = GetCustomValueb("USEWHOLEWINDOWASGLAREA", true)
		this.gtkmenubar = GetCustomValueb("MENUBARFIX", true)
	}

	AssertValue(&VAR_START, "", "StartScene")
	str := generateMain(false)
	file, err := os.Create(WorkingDir() + "main.go")
	if err != nil {
		fmt.Println("Failed to generate main.go:", err)
		os.Exit(1)
	}
	file.WriteString(str)
	file.Close()
}
func (*DesktopBuild) IsGenerated() bool {
	if !FileExists("main.go") {
		return false
	}

	file, err := os.Open("main.go")
	if err != nil {
		return false
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	str := string(contents)

	return (strings.Contains(str, "GLFW") || strings.Contains(str, "SDL") || strings.Contains(str, "GTK")) && strings.Contains(str, "OpenGL")
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
	var vararch string
	var varos string
	if VAR_ARCH == "runtime" {
		vararch = runtime.GOARCH
	} else {
		vararch = VAR_ARCH
	}
	if VAR_OS == "runtime" {
		varos = runtime.GOOS
	} else {
		varos = VAR_OS
	}

	slash := GetSlash()
	exename := PackageName()
	if varos == "windows" {
		exename += ".exe"
	}
	exportpath := "export" + slash + varos
	if varos == "windows" {
		exportpath += slash + vararch
	}
	ExecCommand("mkdir", "-p", exportpath)
	ExecCommand("cp", exename, exportpath+slash+exename)
	ExecCommand("cp", "-r", "assets", exportpath+slash+"assets")
	if VAR_FRAME == "GLFW" {
		if varos == "linux" {
			ExecCommand("mkdir", "-p", exportpath+slash+"lib")
			ExecCommand("cp", "/usr/lib/x86_64-linux-gnu/libopenal.so.1", exportpath+slash+"lib")
		} else if varos == "windows" {
			if vararch == "386" {
				ExecCommand("cp", "C:\\msys64\\mingw32\\bin\\libopenal-1.dll", exportpath)
			} else {
				ExecCommand("cp", "C:\\msys64\\mingw64\\bin\\libopenal-1.dll", exportpath)
			}
		}
	}
}
func (*DesktopBuild) Clean() {
	ExecCommand("go", "clean", "-r", "--cache")
	ExecCommand("rm", "-f", "main.go")
}

func printEnv(forandroid bool) {
	fmt.Println("OS=" + VAR_OS)
	fmt.Println("ARCH=" + VAR_ARCH)
	fmt.Println("FRAME=" + VAR_FRAME)
	fmt.Println("RENDER=" + VAR_RENDER)
	fmt.Println("START=" + VAR_START)
	fmt.Println("CONFIG=" + VAR_CONFIG)
	if forandroid {
		fmt.Println("ANDROID_API=" + VAR_ANDROID_API)
		fmt.Println("ANDROID_KEYSTORE=" + VAR_ANDROID_KEYSTORE)
		fmt.Println("ANDROID_KEYALIAS=" + VAR_ANDROID_KEYALIAS)
		fmt.Println("ANDROID_KEYPWD=" + VAR_ANDROID_KEYPWD)
		fmt.Println("ANDROID_STOREPWD=" + VAR_ANDROID_STOREPWD)
	}
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

func (*DesktopBuild) Env() {
	printEnv(false)
}

func (*AndroidBuild) Build() bool {
	slash := GetSlash()
	ndkHome := os.Getenv("ANDROID_NDK_HOME")
	sysRoot := ndkHome + slash + "platforms" + slash + "android-" + VAR_ANDROID_API + slash + "arch-arm"
	Env = append(Env, []string{
		"CC=arm-linux-androideabi-gcc",
		"CGO_CFLAGS=-O3 -w -D__ANDROID_API__=" + VAR_ANDROID_API + " -I" + ndkHome + "/sysroot/usr/include -I" + ndkHome + "/sysroot/usr/include/arm-linux-androideabi --sysroot=" + sysRoot,
		"CGO_LDFLAGS=-O3 -L" + ndkHome + "/sysroot/usr/lib -L" + ndkHome + "/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/ --sysroot=" + sysRoot,
		"CGO_CXXFLAGS=-O3",
		"CGO_ENABLED=1",
		"GOOS=android",
		"GOARCH=arm",
	}...)

	if err := ExecCommand("go", "build", "-v", "-tags=static", "-buildmode=c-shared", "-ldflags=-s -w -extldflags=-Wl,-soname,libgohome.so", "-o=android/libs/armeabi-v7a/libgohome.so"); err != nil {
		return false
	}

	var assemble string

	if VAR_CONFIG == "DEBUG" {
		assemble = "assembleDebug"
	} else {
		assemble = "assembleRelease"
	}

	if err := ExecCommand("./gradlew", assemble); err != nil {
		return false
	}

	return true
}

func installAPK() bool {
	slash := GetSlash()
	path := "android" + slash + "build" + slash + "outputs" + slash + "apk" + slash
	if VAR_CONFIG == "DEBUG" {
		path += "debug" + slash + "android-debug.apk"
	} else {
		path += "release" + slash + "android-release.apk"
	}

	if err := ExecCommand("adb", "install", "-r", path); err != nil {
		return false
	}
	return true
}

func (this *AndroidBuild) Install() bool {
	if !this.Build() {
		return false
	}

	return installAPK()
}

func doCopy(path string) {
	if err := ExecCommand("cp", "-r", path, WorkingDir()); err != nil {
		fmt.Println("Failed to copy android files")
		os.Exit(1)
	}
}

func setGradleProperties() {
	slash := GetSlash()
	home := os.Getenv("HOME") + slash
	var str string
	file, err := os.Open(home + ".gradle" + slash + "gradle.properties")
	if err == nil {
		contents, _ := ioutil.ReadAll(file)
		cstr := string(contents)
		values := strings.Split(cstr, "\n")
		for _, v := range values {
			if !strings.Contains(v, "=") {
				continue
			}
			keyvalues := strings.Split(v, "=")
			switch keyvalues[0] {
			case "ANDROID_KEYSTORE", "ANDROID_STOREPWD", "ANDROID_KEYALIAS", "ANDROID_KEYPWD":
			default:
				str += keyvalues[0] + "=" + keyvalues[1] + "\n"
			}
		}
	}

	str += "ANDROID_KEYSTORE=" + VAR_ANDROID_KEYSTORE + "\n"
	str += "ANDROID_STOREPWD=" + VAR_ANDROID_STOREPWD + "\n"
	str += "ANDROID_KEYALIAS=" + VAR_ANDROID_KEYALIAS + "\n"
	str += "ANDROID_KEYPWD=" + VAR_ANDROID_KEYPWD + "\n"

	file, err = os.Create(home + ".gradle" + slash + "gradle.properties")
	if err != nil {
		fmt.Println("Failed to create gradle.properties:", err)
		os.Exit(1)
	}
	file.WriteString(str)
	file.Close()
}

func copyAssets() {
	slash := GetSlash()
	ExecCommand("cp", "-r", "assets", "android"+slash+"src"+slash+"main"+slash+"assets"+slash+"assets")
}

func (*AndroidBuild) Generate() {
	if VAR_FRAME != "SDL" {
		fmt.Println("Android is only compatible with SDL")
		VAR_FRAME = "SDL"
	}
	if !strings.Contains(VAR_RENDER, "OpenGLES") {
		fmt.Println("Android is only compatible with OpenGLES")
		fmt.Println("(1) OpenGLES2")
		fmt.Println("(2) OpenGLES3")
		fmt.Println("(3) OpenGLES31")
		fmt.Print("Which version: ")
		version := ConsoleReadi()
		switch version {
		case 1:
			VAR_RENDER = "OpenGLES2"
		case 2:
			VAR_RENDER = "OpenGLES3"
		case 3:
			VAR_RENDER = "OpenGLES31"
		}
	}

	slash := GetSlash()
	gopath := os.Getenv("GOPATH") + slash
	androidpath := gopath + "src" + slash + "github.com" + slash + "PucklaMotzer09" + slash + "GoHomeEngine" + slash + "android" + slash

	doCopy(androidpath + "android")
	doCopy(androidpath + "gradle")
	doCopy(androidpath + "build.gradle")
	doCopy(androidpath + "build_libraries.sh")
	doCopy(androidpath + "gradlew")
	doCopy(androidpath + "gradlew.bat")
	doCopy(androidpath + "settings.gradle")

	appname := GetCustomValue("APPNAME")
	AssertValue(&VAR_ANDROID_API, "", "APILEVEL")
	AssertValue(&VAR_ANDROID_KEYSTORE, "", "KEYSTORE")
	AssertValue(&VAR_ANDROID_KEYALIAS, "", "KEYALIAS")
	AssertValue(&VAR_ANDROID_KEYPWD, "", "KEYPWD")
	AssertValue(&VAR_ANDROID_STOREPWD, "", "STOREPWD")
	AssertValue(&VAR_START, "", "StartScene")
	CustomValues["TITLE"] = appname
	if _, ok := CustomValues["WIDTH"]; !ok {
		CustomValues["WIDTH"] = "1280"
	}
	if _, ok := CustomValues["HEIGHT"]; !ok {
		CustomValues["HEIGHT"] = "720"
	}

	buildgradle := WorkingDir() + "android" + slash + "build.gradle"
	stringsxml := WorkingDir() + "android" + slash + "src" + slash + "main" + slash + "res" + slash + "values" + slash + "strings.xml"

	ReplaceStringinFile(buildgradle, "%APPNAME%", LowerCaseAndNoNumber(appname))
	ReplaceStringinFile(buildgradle, "%APILEVEL%", VAR_ANDROID_API)
	ReplaceStringinFile(stringsxml, "%APPNAME%", appname)

	setGradleProperties()
	copyAssets()

	if VAR_FRAME != "SDL2" {
		fmt.Println("Android is only compatible with SDL2")
		VAR_FRAME = "SDL2"
	}
	if !strings.Contains(VAR_RENDER, "OpenGLES") {
		fmt.Println("Android is only compatible with OpenGLES")
		fmt.Print("Which version (2,3,31): ")
		version := ConsoleRead()
		VAR_RENDER = "OpenGLES" + version
	}

	str := generateMain(true)
	file, err := os.Create(WorkingDir() + "main.go")
	if err != nil {
		fmt.Println("Failed to create main.go:", err)
		os.Exit(1)
	}

	file.WriteString(str)
	file.Close()
}
func (*AndroidBuild) IsGenerated() bool {
	if !FileExists("main.go") || !FileExists("gradlew") {
		return false
	}

	file, err := os.Open("main.go")
	if err != nil {
		return false
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	str := string(contents)

	return strings.Contains(str, "SDL") && strings.Contains(str, "OpenGLES")
}
func (*AndroidBuild) Run() bool {
	if !installAPK() {
		return false
	}

	if err := ExecCommand("adb", "shell", "am", "start", "-n", "com.gohome."+LowerCaseAndNoNumber(CustomValues["APPNAME"])+"/com.example.android.MyGame"); err != nil {
		return false
	}

	return true
}
func (this *AndroidBuild) Export() {
	slash := GetSlash()
	ExecCommand("mkdir", "-p", "export"+slash+"android")
	path := "android" + slash + "build" + slash + "outputs" + slash + "apk" + slash
	if VAR_CONFIG == "DEBUG" {
		path += "debug" + slash + "android-debug.apk"
	} else {
		path += "release" + slash + "android-release.apk"
	}
	ExecCommand("cp", path, "export"+slash+"android"+slash+CustomValues["APPNAME"]+".apk")
}
func (*AndroidBuild) Clean() {
	ExecCommand("rm", "-f", "-r", "android", "gradle", ".gradle", "build.gradle", "gradlew", "gradlew.bat", "settings.gradle", "build_libraries.sh", "main.go")
	ExecCommand("go", "clean", "-r", "--cache")
}
func (*AndroidBuild) Env() {
	printEnv(true)
}

func (this *JSBuild) Build() bool {
	if VAR_CONFIG == "DEBUG" {
		if err := ExecCommand("gopherjs", "build"); err != nil {
			return false
		}
	} else {
		if err := ExecCommand("gopherjs", "build", "-m"); err != nil {
			return false
		}
	}

	return true
}
func (this *JSBuild) Install() bool {
	return this.Build()
}
func (this *JSBuild) Generate() {
	if VAR_FRAME != "JS" {
		fmt.Println("browser is only compatible with JS")
		VAR_FRAME = "JS"
	}

	if VAR_RENDER != "WebGL" {
		fmt.Println("browser is only compatible with WebGL")
		VAR_RENDER = "WebGL"
	}

	AssertValue(&VAR_START, "", "StartScene")
	GetCustomValue("TITLE")
	GetCustomValuei("WIDTH", 640)
	GetCustomValuei("HEIGHT", 480)

	str := generateMain(false)
	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println("Failed to create main.go:", err)
		os.Exit(1)
	}
	file.WriteString(str)
	file.Close()

	createIndexHTML(".")
}
func (this *JSBuild) IsGenerated() bool {
	if !FileExists("main.go") || !FileExists("index.html") {
		return false
	}

	file, err := os.Open("main.go")
	if err != nil {
		return false
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	str := string(contents)

	return strings.Contains(str, "JS") && strings.Contains(str, "WebGL")
}
func (this *JSBuild) Run() bool {
	WriteConfigFile()
	OpenBrowser("http://localhost:8000")
	ExecCommand("python", "-m", "SimpleHTTPServer", "8000")
	return true
}
func (this *JSBuild) Export() {
	slash := GetSlash()
	ExecCommand("mkdir", "-p", "export"+slash+"browser")
	ExecCommand("cp", PackageName()+".js", "export"+slash+"browser")
	ExecCommand("cp", PackageName()+".js.map", "export"+slash+"browser")
	ExecCommand("cp", "-r", "assets", "export"+slash+"browser")
	createIndexHTML("export" + slash + "browser")
}
func (this *JSBuild) Clean() {
	ExecCommand("rm", "main.go", PackageName()+".js", PackageName()+".js.map", "index.html")
	ExecCommand("go", "clean", "-r", "--cache")
}
func (this *JSBuild) Env() {
	printEnv(false)
}
