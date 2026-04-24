package trace

import (
	"context"
	"slices"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	commonattribute "go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
	otelattributeadapter "go.mws.cloud/go-sdk/internal/client/interceptors/attribute/otel"
	"go.mws.cloud/go-sdk/internal/client/interceptors/retry"
)

const attemptAttribute = attribute.Key("req.attempt")

func New(options ...Option) commonclient.Interceptor {
	cfg := defaultConfig()
	for _, option := range options {
		option.apply(cfg)
	}

	if cfg.tracer == nil {
		cfg.tracer = otel.GetTracerProvider().Tracer(cfg.name)
	}

	i := &traceInvoker{
		tracer:     cfg.tracer,
		attributes: cfg.attributes,
	}
	return i.interceptor
}

type traceInvoker struct {
	tracer     trace.Tracer
	attributes []attribute.KeyValue
}

func (i *traceInvoker) interceptor(ctx context.Context, request any, response commonclient.APIResp, invoker commonclient.Invoker) error {
	ctx = i.checkContext(ctx)
	attributes := append(otelattributeadapter.Convert(commonattribute.FromContext(ctx)), i.attributes...)
	attributes = append(attributes, attemptAttribute.Int(retry.AttemptFromContext(ctx)))

	ctx, span := i.tracer.Start(ctx,
		i.spanName(attributes),
		trace.WithAttributes(attributes...),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	err := invoker(ctx, request, response)
	if err == nil {
		span.SetAttributes(i.extractResponseAttributes(response)...)
		return nil
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	return err
}

func (i *traceInvoker) extractResponseAttributes(response commonclient.APIResp) []attribute.KeyValue {
	if response == nil {
		return nil
	}

	return []attribute.KeyValue{attribute.Int("code", response.GetCode())}
}

func (i *traceInvoker) checkContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

func (i *traceInvoker) spanName(attributes []attribute.KeyValue) string {
	name := "cloud.mws.api"

	if projectName := i.getProjectName(attributes); projectName != "" {
		name += "." + projectName
	}
	if serviceName := i.getServiceName(attributes); serviceName != "" {
		name += "." + serviceName
	}
	name += ".client"
	if method := i.getMethod(attributes); method != "" {
		name += "." + method
	}

	return name
}

func (i *traceInvoker) getProjectName(attributes []attribute.KeyValue) string {
	return i.getAttribute(attributes, commonattribute.ProjectName)
}

func (i *traceInvoker) getServiceName(attributes []attribute.KeyValue) string {
	return i.getAttribute(attributes, commonattribute.ServiceName)
}

func (i *traceInvoker) getMethod(attributes []attribute.KeyValue) string {
	return i.getAttribute(attributes, commonattribute.Method)
}

func (i *traceInvoker) getAttribute(attributes []attribute.KeyValue, attributeName attribute.Key) string {
	index := slices.IndexFunc(attributes, func(value attribute.KeyValue) bool {
		return value.Key == attributeName
	})

	if index == -1 {
		return ""
	}

	return attributes[index].Value.AsString()
}
