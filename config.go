package httpsrv

import "io"

type config struct {
	port int
	path string
	out  io.Writer
}

type Options func(cfg *config)

func WithPort(port int) Options {
	return func(cfg *config) {
		cfg.port = port
	}
}

func WithPath(path string) Options {
	return func(cfg *config) {
		cfg.path = path
	}
}

func WithLogger(out io.Writer) Options {
	return func(cfg *config) {
		cfg.out = out
	}
}

func buildConfig(opts ...Options) *config {
	cfg := &config{
		path: ".",
	}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}
