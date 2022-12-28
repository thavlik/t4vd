package api

import "github.com/thavlik/bjjvb/base/pkg/base"

func NewFilterClientFromOptions(opts base.ServiceOptions) Filter {
	options := NewFilterClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewFilterClient(opts.Endpoint, options)
}
