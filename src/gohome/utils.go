package gohome

import (
	"io"
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
