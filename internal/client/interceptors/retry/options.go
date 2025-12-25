package retry

// Option applies a configuration to the given config
type Option interface {
	apply(cfg *config)
}

type optionFunc func(cfg *config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

// WithRetryer sets the retryer. By default, NoRetry is used
func WithRetryer(retryer Retryer) Option {
	return optionFunc(func(cfg *config) {
		cfg.retryer = retryer
	})
}

func WithFailOnExhaustedAttempts(fail bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.failOnExhaustedAttempts = fail
	})
}
