package http

import (
	"net/http"
	"time"

	"github.com/everfir/go-helpers/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func NewTraceTripper(transport http.RoundTripper) *TraceTripper {
	return &TraceTripper{transport: transport}
}

// TraceTripper is a middleware that wraps the default transport
type TraceTripper struct {
	transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface to add middleware logic
func (c *TraceTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Middleware logic before sending the request
	// Get the context from the request
	ctx := req.Context()
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// header
	req.Header.Add(env.BusinessKey.String(), env.Business(ctx))
	req.Header.Add(env.VersionKey.String(), env.Version(ctx))
	req.Header.Add(env.PlatformKey.String(), env.Platform(ctx).String())
	req.Header.Add(env.DeviceKey.String(), env.Device(ctx).String())
	req.Header.Add(env.AppTypeKey.String(), env.AppType(ctx).String())

	// Call the next RoundTripper (default transport in this case)
	resp, err := c.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var (
	Client = http.Client{
		Timeout: 10 * time.Second,
		Transport: &TraceTripper{
			transport: http.DefaultTransport,
		},
	}
)
