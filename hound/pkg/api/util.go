package api

import "github.com/thavlik/t4vd/base/pkg/base"

func NewHoundClientFromOptions(opts base.ServiceOptions) Hound {
	options := NewHoundClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewHoundClient(opts.Endpoint, options)
}
