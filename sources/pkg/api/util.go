package api

import "github.com/thavlik/t4vd/base/pkg/base"

func NewSourcesClientFromOptions(opts base.ServiceOptions) Sources {
	options := NewSourcesClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewSourcesClient(opts.Endpoint, options)
}
