package cachedset

import "context"

type CachedSet interface {
	Set(ctx context.Context, key string, value string, index int) error
	Complete(ctx context.Context, key string) error
	List(ctx context.Context, key string) (values []string, complete bool, err error)
}
