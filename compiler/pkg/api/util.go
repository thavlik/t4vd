package api

import "github.com/thavlik/t4vd/base/pkg/base"

func NewCompilerClientFromOptions(opts base.ServiceOptions) Compiler {
	options := NewCompilerClientOptions().SetTimeout(opts.Timeout)
	if opts.BasicAuth.Username != "" {
		options.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	return NewCompilerClient(opts.Endpoint, options)
}
