package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getGradleProperties() (data map[string]string) {
	file := openGradleProperties()
	if file != nil {
		data = make(map[string]string)
		contents, _ := ioutil.ReadAll(file)
		cstr := string(contents)
		values := strings.Split(cstr, "\n")
		for _, v := range values {
			if !strings.Contains(v, "=") {
				continue
			}
			keyvalues := strings.Split(v, "=")
			data[keyvalues[0]] = keyvalues[1]
		}
		file.Close()
	}
	return
}

func openGradleProperties() *os.File {
	slash := GetSlash()
	home := os.Getenv("HOME") + slash
	file, err := os.Open(home + ".gradle" + slash + "gradle.properties")
	if err == nil {
		return file
	}
	return nil
}

func writeGradleProperties(file *os.File, data map[string]string) {
	for k, v := range data {
		file.WriteString(k + "=" + v + "\n")
	}
}

func setGradleProperties() {
	props := getGradleProperties()

	if props == nil {
		props = make(map[string]string)
	}

	props["ANDROID_KEYSTORE"] = VAR_ANDROID_KEYSTORE
	props["ANDROID_KEYPWD"] = VAR_ANDROID_KEYPWD
	props["ANDROID_STOREPWD"] = VAR_ANDROID_STOREPWD
	props["ANDROID_KEYALIAS"] = VAR_ANDROID_KEYALIAS

	file, err := os.Create(os.Getenv("HOME") + GetSlash() + ".gradle" + GetSlash() + "gradle.properties")
	if err != nil {
		fmt.Println("Error create gradle.properties:", err)
		os.Exit(1)
	}
	writeGradleProperties(file, props)
	file.Close()
}

func setExistingGradleProperties() {
	props := getGradleProperties()
	if props != nil {
		if val, ok := props["ANDROID_KEYSTORE"]; ok {
			VAR_ANDROID_KEYSTORE = val
		}
		if val, ok := props["ANDROID_KEYPWD"]; ok {
			VAR_ANDROID_KEYPWD = val
		}
		if val, ok := props["ANDROID_STOREPWD"]; ok {
			VAR_ANDROID_STOREPWD = val
		}
		if val, ok := props["ANDROID_KEYALIAS"]; ok {
			VAR_ANDROID_KEYALIAS = val
		}
	}
}
