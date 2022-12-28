package slideshow

import (
	"fmt"
	"io"
)

type readAhead struct {
	r      io.ReadSeeker //
	rc     int64
	bufCur int64  // the offset that the buffer starts at
	bufLen int    // valid length of the buffer
	buf    []byte // the buffer itself
	c      int64  // desired cursor
}

func NewReadAhead(r io.ReadSeeker, bufCap int) io.ReadSeeker {
	return &readAhead{
		r:      r,
		buf:    make([]byte, bufCap, bufCap),
		bufCur: -1,
	}
}

func (a *readAhead) isCursorInBuf() bool {
	return a.bufCur <= a.c && a.c < a.bufCur+int64(a.bufLen)
}

func (a *readAhead) canReadFromBuf() bool {
	return a.isValidBuf() && a.isCursorInBuf()
}

func (a *readAhead) isValidBuf() bool { return a.bufCur != -1 }

func (a *readAhead) invalidateBuf() { a.bufCur = -1 }

func (a *readAhead) readFromBuf(p []byte) (n int) {
	// The start of where we're supposed to read is within the buffer.
	n = copy(p, a.buf[a.c-a.bufCur:a.bufLen])
	a.c += int64(n) // new cursor position
	return
}

func (a *readAhead) readUnderlying(p []byte) (n int, err error) {
	a.bufLen, err = a.r.Read(a.buf) // read underlying to buf
	if err != nil {
		return 0, err
	}
	a.bufCur = a.rc // set new buf cursor
	a.rc += int64(a.bufLen)
	if actualRc, _ := a.Seek(0, io.SeekCurrent); actualRc != a.rc {
		panic(fmt.Sprintf("rc is %d, expected %d", a.rc, actualRc))
	}
	return a.readFromBuf(p), nil
}

func (a *readAhead) Read(p []byte) (n int, err error) {
	if a.canReadFromBuf() {
		// The start of where we're supposed to read is within the buffer.
		n = a.readFromBuf(p)
		if n == len(p) {
			// Everything was read from the buffer
			return
		}
		p = p[n:]
	}

	if _, err := a.Seek(a.c, io.SeekStart); err != nil {
		return n, err
	}

	// The cursors should always match, even if the buffer is invalid.
	// If the buffer is invalid because of a seek operation, then
	// the cursor should match because it was set at the same time.
	if got, expected := a.c, a.rc; got != expected {
		panic(fmt.Sprintf("cursor sanity check failed (got %d, expected %d, cursor is ahead by %d)", got, expected, got-expected))
	}

	if len(p) >= len(a.buf) {
		// Do a direct read. Buffer does not provide
		// a performance advantage here.
		n2, err := a.r.Read(p)
		if err != nil {
			return 0, err
		}
		a.c += int64(n2)
		a.rc += int64(n2)
		return n + n2, nil
	}

	// The remaining bytes will require reading from the
	// underlying stream.
	n2, err := a.readUnderlying(p)
	if err != nil {
		return 0, err
	}
	return n + n2, nil
}

func (a *readAhead) Seek(offset int64, whence int) (n int64, err error) {
	if whence == io.SeekCurrent {
		n, err = a.r.Seek(a.c+offset, io.SeekStart)
		if err != nil {
			return
		}
	} else {
		n, err = a.r.Seek(offset, whence)
		if err != nil {
			return
		}
	}
	a.rc = n
	a.c = n
	return
}
