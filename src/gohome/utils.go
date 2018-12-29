package gohome

import (
	"io"
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
	const NUM_FLOATS = MESH2DVERTEX_SIZE / 4
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
