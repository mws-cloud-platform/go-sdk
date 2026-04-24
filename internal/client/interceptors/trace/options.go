package trace

import (
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Option applies a configuration to the given config
type Option interface {
	apply(cfg *config)
}

type optionFunc func(cfg *config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

// WithName sets the prefix for the spans created by this invoker.
// by default, the "mws.client" name is used.
func WithName(name string) Option {
	return optionFunc(func(cfg *config) {
		if name = strings.TrimSpace(name); name != "" {
			cfg.name = DefaultName + "." + name
		}
	})
}

// WithTracer sets the underlying trace.Tracer that is used to create spans
// By default the no-op tracer is used.
func WithTracer(tracer trace.Tracer) Option {
	return optionFunc(func(cfg *config) {
		cfg.tracer = tracer
	})
}

// WithAttributes sets a list of attribute.KeyValue labels for all spans associated with this invoker
func WithAttributes(attributes ...attribute.KeyValue) Option {
	return optionFunc(func(cfg *config) {
		cfg.attributes = attributes
	})
}
