package main

var (
	// Common
	VAR_OS     string = "runtime"
	VAR_ARCH   string = "runtime"
	VAR_FRAME  string = "GLFW"
	VAR_RENDER string = "OpenGL"
	VAR_START  string = ""
	VAR_CONFIG string = "DEBUG"

	// Android
	VAR_ANDROID_API      string = ""
	VAR_ANDROID_KEYSTORE string = ""
	VAR_ANDROID_KEYALIAS string = ""
	VAR_ANDROID_KEYPWD   string = ""
	VAR_ANDROID_STOREPWD string = ""

	COMMAND string = ""

	Env []string

	CustomValues map[string]string

	old_OS               string = ""
	old_ARCH             string = ""
	old_FRAME            string = ""
	old_RENDER           string = ""
	old_START            string = ""
	old_CONFIG           string = ""
	old_ANDROID_API      string = ""
	old_ANDROID_KEYSTORE string = ""
	old_ANDROID_KEYALIAS string = ""
	old_ANDROID_KEYPWD   string = ""
	old_ANDROID_STOREPWD string = ""
	old_CustomValues     map[string]string
)
