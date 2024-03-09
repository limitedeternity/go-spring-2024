//go:build !solution

package otp

import (
	"io"
)

type CipherReader struct {
	r    io.Reader
	prng io.Reader
}

func preserveError[T any](v T, err error) func(error) (T, error) {
	return func(prev error) (T, error) {
		if prev != nil {
			return v, prev
		}

		return v, err
	}
}

func (c *CipherReader) Read(p []byte) (n int, err error) {
	n, err = preserveError(c.r.Read(p))(err)
	rng := make([]byte, n)
	_, err = preserveError(c.prng.Read(rng))(err)

	for i := 0; i < n; i++ {
		p[i] ^= rng[i]
	}

	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &CipherReader{r, prng}
}

type CipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (c *CipherWriter) Write(p []byte) (n int, err error) {
	rng := make([]byte, len(p))
	_, err = preserveError(c.prng.Read(rng))(err)

	for i := 0; i < len(p); i++ {
		rng[i] ^= p[i]
	}

	n, err = preserveError(c.w.Write(rng))(err)
	return
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &CipherWriter{w, prng}
}
