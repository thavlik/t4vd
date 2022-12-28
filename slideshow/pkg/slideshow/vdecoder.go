package slideshow

import (
	"fmt"
	"image"
	"time"

	"github.com/ebml-go/webm"
	"github.com/pkg/errors"
	"github.com/xlab/libvpx-go/vpx"
)

type Frame struct {
	image.Image
	Timecode   time.Duration
	IsKeyframe bool
}

type VDecoder struct {
	src   <-chan webm.Packet
	ctx   *vpx.CodecCtx
	iface *vpx.CodecIface
}

type VCodec string

const (
	CodecVP8  VCodec = "V_VP8"
	CodecVP9  VCodec = "V_VP9"
	CodecVP10 VCodec = "V_VP10"
)

func NewVDecoder(codec VCodec, src <-chan webm.Packet) (*VDecoder, error) {
	dec := &VDecoder{
		src: src,
		ctx: vpx.NewCodecCtx(),
	}
	switch codec {
	case CodecVP8:
		dec.iface = vpx.DecoderIfaceVP8()
	case CodecVP9:
		dec.iface = vpx.DecoderIfaceVP9()
	default:
		return nil, errors.New("unsupported codec")
	}
	if err := vpx.Error(vpx.CodecDecInitVer(
		dec.ctx,
		dec.iface,
		nil,
		0,
		vpx.DecoderABIVersion,
	)); err != nil {
		return nil, errors.Wrap(err, "CodecDecInitVer")
	}
	return dec, nil
}

func (v *VDecoder) Process(
	out chan<- Frame,
	stop <-chan struct{},
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
		close(out)
	}()
	for pkt := range v.src {
		dataSize := uint32(len(pkt.Data))
		if err := vpx.Error(vpx.CodecDecode(
			v.ctx,
			string(pkt.Data),
			dataSize,
			nil,
			0,
		)); err != nil {
			return errors.Wrap(err, "CodecDecode")
		}
		var iter vpx.CodecIter
		img := vpx.CodecGetFrame(v.ctx, &iter)
		for img != nil {
			img.Deref()
			frame := Frame{
				Image:    img.ImageYCbCr(),
				Timecode: pkt.Timecode,
			}
			select {
			case <-stop:
				return nil
			case out <- frame:
				img = vpx.CodecGetFrame(v.ctx, &iter)
				continue
			}
		}
	}
	return nil
}
