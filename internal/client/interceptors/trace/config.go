package trace

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultName = "mws.client"
)

type config struct {
	name       string
	tracer     trace.Tracer
	attributes []attribute.KeyValue
}

func defaultConfig() *config {
	return &config{
		name: DefaultName,
	}
}
