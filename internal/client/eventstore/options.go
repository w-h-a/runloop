package eventstore

import "context"

type Option func(*Options)

type Options struct {
	Location string
	Context  context.Context
}

func WithLocation(location string) Option {
	return func(o *Options) {
		o.Location = location
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
