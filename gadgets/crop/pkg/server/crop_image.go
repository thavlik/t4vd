package server

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/pkg/errors"
)

func cropImage(
	box image.Rectangle,
	r io.Reader,
	w io.Writer,
) error {
	img, ty, err := image.Decode(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode image")
	}
	img, err = doCropImage(img, box)
	if err != nil {
		return errors.Wrap(err, "failed to crop image")
	}
	switch ty {
	case "jpeg":
		if err := jpeg.Encode(
			w,
			img,
			&jpeg.Options{Quality: 100},
		); err != nil {
			return errors.Wrap(err, "failed to encode jpeg")
		}
		return nil
	case "png":
		if err := png.Encode(
			w,
			img,
		); err != nil {
			return errors.Wrap(err, "failed to encode png")
		}
		return nil
	default:
		return errors.Errorf("unsupported image type: %s", ty)
	}
}

func doCropImage(
	img image.Image,
	crop image.Rectangle,
) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}
