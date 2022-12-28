package slideshow

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
)

type s3ReadSeeker struct {
	bucket    *string
	key       *string
	s3        *s3.S3
	c         int64
	l         int64
	totalRead int64
}

var ErrNullObject = errors.New("object has zero length")

var chunkSize = 1024

func NewS3ReadSeeker(bucket, key string) (io.ReadSeeker, error) {
	b := aws.String(bucket)
	k := aws.String(key)
	s3Client := s3.New(base.AWSSession())
	resp, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: b,
		Key:    k,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s3")
	}
	l := aws.Int64Value(resp.ContentLength)
	if l == 0 {
		return nil, ErrNullObject
	}
	return &s3ReadSeeker{
		bucket: b,
		key:    k,
		s3:     s3Client,
		l:      l,
	}, nil
}

func (h *s3ReadSeeker) Read(p []byte) (n int, err error) {
	start := h.c
	end := start + int64(len(p))
	if end > h.l {
		end = h.l
	}
	rn := fmt.Sprintf(
		"bytes=%d-%d",
		start,
		end-1,
	)
	result, err := h.s3.GetObject(&s3.GetObjectInput{
		Bucket: h.bucket,
		Key:    h.key,
		Range:  aws.String(rn),
	})
	if err != nil {
		return 0, errors.Wrap(err, "s3")
	}
	defer result.Body.Close()
	b := make([]byte, len(p))
	n, err = result.Body.Read(b)
	for i, v := range b {
		p[i] = v
	}
	h.c += int64(n)
	h.totalRead += int64(n)
	return n, err
}

func (h *s3ReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		h.c = offset
	case io.SeekCurrent:
		h.c += offset
	case io.SeekEnd:
		h.c = h.l - offset
	default:
		panic(base.Unreachable)
	}
	return h.c, nil
}
