package expandenv

// Option is config option.
type Option func(*options)

type options struct {
	getenv func(string) string
}

func WithGetenv(getenv func(string) string) Option {
	return func(opts *options) {
		opts.getenv = getenv
	}
}
