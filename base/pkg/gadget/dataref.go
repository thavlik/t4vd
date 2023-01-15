package gadget

import (
	"context"

	"github.com/pkg/errors"
)

var ErrNullDataRef = errors.New("null data ref")

type DataRef struct {
	gadget  string
	channel string
	l       chan struct{}
}

func NewDataRef(
	gadget string,
	channel string,
) *DataRef {
	return &DataRef{
		gadget,
		channel,
		make(chan struct{}, 1),
	}
}

func (g *DataRef) Get(ctx context.Context) (gadget string, channel string, err error) {
	if err := g.lock(ctx); err != nil {
		return "", "", errors.Wrap(err, "lock")
	}
	defer g.unlock()
	if g.gadget == "" || g.channel == "" {
		return "", "", ErrNullDataRef
	}
	return g.gadget, g.channel, nil
}

func (g *DataRef) Set(
	ctx context.Context,
	endpoint string,
	channel string,
) error {
	if err := g.lock(ctx); err != nil {
		return errors.Wrap(err, "lock")
	}
	defer g.unlock()
	g.gadget = endpoint
	g.channel = channel
	return nil
}

func (g *DataRef) lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case g.l <- struct{}{}:
		return nil
	}
}

func (g *DataRef) unlock() {
	select {
	case <-g.l:
	default:
		panic("unlocking unlocked gadget ref")
	}
}
