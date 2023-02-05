package tracing

type (
	Option func(o *tracingOptions)

	tracingOptions struct {
		url         string
		serviceName string
		envName     string
	}
)

func WithURL(url string) Option {
	return func(o *tracingOptions) {
		o.url = url
	}
}

func WithService(service string) Option {
	return func(o *tracingOptions) {
		o.serviceName = service
	}
}

func WithEnvironment(env string) Option {
	return func(o *tracingOptions) {
		o.envName = env
	}
}

func evaluateOptions(opts ...Option) tracingOptions {
	options := tracingOptions{}

	for _, optFunc := range opts {
		optFunc(&options)
	}

	return options
}
