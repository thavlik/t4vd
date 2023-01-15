package server

import (
	"image"
	"io"

	"github.com/pkg/errors"
)

func cropVideo(
	box image.Rectangle,
	marker *Marker,
	r io.Reader,
	w io.Writer,
) error {
	return errors.New("not implemented")
}
