//go:build exclude

package solutions

import "golang.org/x/tour/reader"

type MyReader struct{}

func (r MyReader) Read(b []byte) (n int, err error) {
	n, err = len(b), nil

	for i := 0; i < n; i++ {
		b[i] = 'A'
	}

	return
}

func main() {
	reader.Validate(MyReader{})
}
