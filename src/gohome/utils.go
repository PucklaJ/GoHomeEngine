package gohome

import (
	"io"
	"strings"
	"sync"
)

const (
	READ_ALL_BUFFER_SIZE = 512 * 512
)

func ReadAll(r io.Reader) (str string, err error) {
	str = ""
	var n int = 1
	for err == nil && n != 0 {
		buf := make([]byte, READ_ALL_BUFFER_SIZE)
		n, err = r.Read(buf)
		str += string(buf[:n])
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func Maxi(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Mini(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Mesh3DVerticesToFloatArray(vertices []Mesh3DVertex) (array []float32) {
	const NUM_FLOATS = MESH3DVERTEXSIZE / 4
	array = make([]float32, len(vertices)*NUM_FLOATS)
	var wg sync.WaitGroup
	wg.Add(len(array) / NUM_FLOATS)
	var index = 0
	for _, v := range vertices {
		go func(_index int, _v Mesh3DVertex) {
			for i := 0; i < NUM_FLOATS; i++ {
				array[_index+i] = _v[i]
			}
			wg.Done()
		}(index, v)
		index += NUM_FLOATS
	}
	wg.Wait()
	return
}

func Mesh2DVerticesToFloatArray(vertices []Mesh2DVertex) (array []float32) {
	const NUM_FLOATS = MESH2DVERTEXSIZE / 4
	array = make([]float32, len(vertices)*NUM_FLOATS)
	var wg sync.WaitGroup
	wg.Add(len(array) / NUM_FLOATS)
	var index = 0
	for _, v := range vertices {
		go func(_index int, _v Mesh2DVertex) {
			for i := 0; i < NUM_FLOATS; i++ {
				array[_index+i] = _v[i]
			}
			wg.Done()
		}(index, v)
		index += NUM_FLOATS
	}
	wg.Wait()
	return
}

func Lines3DToFloatArray(lines []Line3D) (array []float32) {
	const NUM_FLOATS = SHAPE3DVERTEXSIZE / 4 * 2
	array = make([]float32, len(lines)*NUM_FLOATS)
	var wg sync.WaitGroup
	wg.Add(len(array) / NUM_FLOATS)
	var index = 0
	for _, l := range lines {
		go func(_index int, _l Line3D) {
			for i := 0; i < 2; i++ {
				for j := 0; j < 3+4; j++ {
					array[_index+i*(3+4)+j] = _l[i][j]
				}
			}
			wg.Done()
		}(index, l)
		index += NUM_FLOATS
	}
	wg.Wait()
	return
}

func Shape2DVerticesToFloatArray(vertices []Shape2DVertex) (array []float32) {
	const NUM_FLOATS = SHAPE2DVERTEXSIZE / 4
	array = make([]float32, len(vertices)*NUM_FLOATS)
	var wg sync.WaitGroup
	wg.Add(len(array) / NUM_FLOATS)
	var index = 0
	for _, v := range vertices {
		go func(_index int, _v Shape2DVertex) {
			for i := 0; i < NUM_FLOATS; i++ {
				array[_index+i] = _v[i]
			}
			wg.Done()
		}(index, v)
		index += NUM_FLOATS
	}
	wg.Wait()
	return
}

func EqualIgnoreCase(str1, str string) bool {
	if len(str1) != len(str) {
		return false
	}
	for i := 0; i < len(str1); i++ {
		if str1[i] != str[i] {
			if str1[i] >= 65 && str1[i] <= 90 {
				if str[i] >= 97 && str[i] <= 122 {
					if str1[i]+32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else if str1[i] >= 97 && str1[i] <= 122 {
				if str[i] >= 65 && str[i] <= 90 {
					if str1[i]-32 != str[i] {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func GetFileExtension(file string) string {
	index := strings.LastIndex(file, ".")
	if index == -1 {
		return ""
	}
	return file[index+1:]
}

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

func OpenFileWithPaths(path string, paths []string) (File, string, error) {
	var reader File
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
