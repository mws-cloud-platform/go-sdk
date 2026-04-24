package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestWithName(t *testing.T) {
	cfg := defaultConfig()
	WithName("  test").apply(cfg)
	assert.Equal(t, DefaultName+".test", cfg.name)
}

func TestWithAttributes(t *testing.T) {
	cfg := defaultConfig()
	at := []attribute.KeyValue{
		attribute.String("s", "v"),
		attribute.Bool("b", true),
	}
	WithAttributes(at...).apply(cfg)
	assert.Equal(t, at, cfg.attributes)
}

func TestWithTracer(t *testing.T) {
	cfg := defaultConfig()
	tracer := noop.NewTracerProvider().Tracer(DefaultName)
	WithTracer(tracer).apply(cfg)
	assert.Equal(t, tracer, cfg.tracer)
}
