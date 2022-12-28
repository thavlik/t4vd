package base

import "context"

func Progress(ctx context.Context, onProgress chan<- struct{}) {
	if onProgress != nil {
		select {
		case <-ctx.Done():
			return
		case onProgress <- struct{}{}:
		}
	}
}
