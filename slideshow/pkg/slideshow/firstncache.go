package slideshow

import (
	"io"

	"github.com/pkg/errors"
)

type firstNCache struct {
	r   io.ReadSeeker //
	buf []byte        // cache of the first n bytes
	c   int64         // desired cursor
}

func NewFirstNCache(r io.ReadSeeker, bufSize int) (io.ReadSeeker, error) {
	a := &firstNCache{
		r:   r,
		buf: make([]byte, bufSize, bufSize),
	}
	n, err := a.r.Read(a.buf)
	if err != nil {
		return nil, errors.Wrap(err, "read")
	}
	a.buf = a.buf[:n]
	return a, nil
}

func (a *firstNCache) Read(p []byte) (n int, err error) {
	bufLen := int64(len(a.buf))
	if a.c < bufLen {
		// At least some of the bytes are cached
		n = copy(p, a.buf[a.c:])
		a.c += int64(n)
		// Check if we have more bytes to read
		if n == len(p) {
			// No more bytes to read
			return n, nil
		}
		p = p[n:]
		/*
			// Read the remainder from the underlying reader
			if _, err := a.r.Seek(bufLen, 0); err != nil {
				return 0, errors.Wrap(err, "seek")
			}
			n2, err := a.r.Read(p[n:])
			if err != nil {
				return 0, errors.Wrap(err, "read")
			}
			a.c += int64(n2)

			// Consider both sources when calculating num bytes written
			return n + n2, nil
		*/
	}
	n, err = a.r.Read(p)
	a.c += int64(n)
	return n, err
}

func (a *firstNCache) Seek(offset int64, whence int) (n int64, err error) {
	n, err = a.r.Seek(offset, whence)
	a.c = n
	return n, err
}
