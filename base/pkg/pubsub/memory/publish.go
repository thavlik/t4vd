package redis

func (p *memoryPubSub) Publish(payload []byte) error {
	p.l.Lock()
	defer p.l.Unlock()
	for ch := range p.channels {
		select {
		case ch <- payload:
		default:
			p.log.Warn("memory pubsub dropped message due to channel being full")
		}
	}
	return nil
}
