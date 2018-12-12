package gohome

import (
	"strings"
)

func GetFileFromPath(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[index+1:]
	} else {
		return path
	}
}

func GetPathFromFile(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[:index+1]
	} else {
		return ""
	}
}

func OpenFileWithPaths(path string, paths []string) (*File, string, error) {
	var reader *File
	var err error
	var filename string

	for i := 0; i < len(paths); i++ {
		filename = paths[i] + path
		if reader, err = Framew.OpenFile(filename); err == nil {
			break
		} else if reader, err = Framew.OpenFile(paths[i] + GetFileFromPath(path)); err == nil {
			filename = paths[i] + GetFileFromPath(path)
			break
		}
	}

	return reader, filename, err
}
