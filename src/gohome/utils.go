package gohome

import (
	"io"
	"strings"
	"sync"
)

const (
	READ_ALL_BUFFER_SIZE = 512 * 512
)

// Reads the entire content of a reader. Uses a bigger buffer than the normal one
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

// Returns the maximum value of a and b
func Maxi(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

// Returns the minimum value of a and b
func Mini(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Converts a mesh 3d vertex array to a float array used by OpenGL
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

// Converts a mesh 2d vertex array to a float array used by OpenGL
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

// Converts a shape 3d vertex array to a float array used by OpenGL
func Shape3DVerticesToFloatArray(points []Shape3DVertex) (array []float32) {
	const NUM_FLOATS = SHAPE3DVERTEXSIZE / 4
	array = make([]float32, len(points)*NUM_FLOATS)
	var wg sync.WaitGroup
	wg.Add(len(array) / NUM_FLOATS)
	var index = 0
	for _, p := range points {
		go func(_index int, _p Shape3DVertex) {
			for j := 0; j < 3+4; j++ {
				array[_index+j] = _p[j]
			}
			wg.Done()
		}(index, p)
		index += NUM_FLOATS
	}
	wg.Wait()
	return
}

// Converts a shape 2d vertex array to a float array used by OpenGL
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

// Returns wether one string equals the other ignoring the case
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

// The returns the file extension of a file name
func GetFileExtension(file string) string {
	index := strings.LastIndex(file, ".")
	if index == -1 {
		return ""
	}
	return file[index+1:]
}

// Returns the file name of a file path
func GetFileFromPath(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[index+1:]
	} else {
		return path
	}
}

// Returns the directory of a file
func GetPathFromFile(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[:index+1]
	} else {
		return ""
	}
}

// Opens a file in multiple paths returning the first one that works
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
