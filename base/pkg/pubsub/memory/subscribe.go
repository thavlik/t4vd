package redis

import "context"

func (p *memoryPubSub) Subscribe(ctx context.Context) (<-chan []byte, error) {
	p.l.Lock()
	defer p.l.Unlock()
	ch := make(chan []byte, 32)
	p.channels[ch] = struct{}{}
	return ch, nil
}
