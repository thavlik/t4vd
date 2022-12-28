package slideshow

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ebml-go/webm"
	"github.com/pkg/errors"
)

var defaultReadAhead = 8192 // 8kb

// GetFrameFromReader reads a webm file, seeks to the time,
// and returns the next decoded frame.
func GetSingleFrameFromReader(
	r io.ReadSeeker,
	t time.Duration,
) (*Frame, error) {
	frames := make(chan Frame, 1)
	stop := make(chan struct{}, 1)
	done := make(chan error, 1)
	go func() {
		done <- GetFramesFromReader(r, t, frames, stop)
	}()
	frame, ok := <-frames
	stop <- struct{}{}
	if !ok {
		return nil, errors.New("no frames decoded before channel closed")
	}
	if err := <-done; err != nil {
		return nil, errors.Wrap(err, "GetFramesFromReader")
	}
	return &frame, nil
}

func GetSingleFrameFromBucket(
	bucket string,
	key string,
	t time.Duration,
) (*Frame, error) {
	s3r, err := NewS3ReadSeeker(
		bucket,
		key,
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewS3ReadSeeker")
	}
	return GetSingleFrameFromReader(
		NewReadAhead(s3r, defaultReadAhead),
		t,
	)
}

func GetSingleFrameFromFile(
	path string,
	t time.Duration,
) (*Frame, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}
	defer f.Close()
	return GetSingleFrameFromReader(f, t)
}

func GetFramesFromReader(
	r io.ReadSeeker,
	t time.Duration,
	frame chan<- Frame,
	stop <-chan struct{},
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	hdr := &webm.WebM{}
	pr, err := webm.Parse(r, hdr)
	if err != nil {
		return errors.Wrap(err, "webm.Parse")
	}
	pr.Seek(t)
	vtrack := hdr.FindFirstVideoTrack()
	if vtrack == nil {
		return errors.New("missing video track")
	}
	vPackets := make(chan webm.Packet, 32)
	vdec, err := NewVDecoder(VCodec(vtrack.CodecID), vPackets)
	if err != nil {
		return errors.Wrap(err, "NewVDecoder")
	}
	stopDec := make(chan struct{}, 1)
	d0 := make(chan struct{}, 1)
	d1 := make(chan error, 1)
	go func() { // demuxer
		defer func() {
			stopDec <- struct{}{}
			close(vPackets)
			pr.Shutdown()
			d0 <- struct{}{}
		}()
		for {
			select {
			case <-stop:
				return
			case pkt, ok := <-pr.Chan:
				if !ok {
					return
				}
				if pkt.TrackNumber == vtrack.TrackNumber {
					select {
					case <-stop:
						return
					case vPackets <- pkt:
					}
				}
			}
		}
	}()
	go func() {
		d1 <- vdec.Process(frame, stopDec)
	}()
	if _, err = <-d0, <-d1; err != nil {
		return errors.Wrap(err, "Process")
	}
	return nil
}

func GetFramesFromFile(
	path string,
	t time.Duration,
	frame chan<- Frame,
	stop <-chan struct{},
) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()
	return GetFramesFromReader(f, t, frame, stop)
}

func GetFramesFromBucket(
	bucket string,
	key string,
	t time.Duration,
	frame chan<- Frame,
	stop <-chan struct{},
) error {
	s3r, err := NewS3ReadSeeker(
		bucket,
		key,
	)
	if err != nil {
		return errors.Wrap(err, "NewS3ReadSeeker")
	}
	return GetFramesFromReader(
		NewReadAhead(s3r, defaultReadAhead),
		t,
		frame,
		stop,
	)
}
