package etcd

import "strings"

type Option func(*options)

type options struct {
	root  string
	paths []string
}

func WithRoot(prefix string) Option {
	if !strings.HasSuffix(prefix, "/") { // 维护路径
		prefix += "/"
	}
	return func(o *options) {
		o.root = prefix
	}
}

func WithPaths(path ...string) Option {
	return func(o *options) {
		o.paths = path
	}
}

func NewOptions(opts ...Option) *options {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
