package strata

import "context"

type metricsKey struct{}

// FromContext extracts and returns the Metrics from the context.  An error is
// returned if the context does not contain Metrics or the context is nil.
func FromContext(ctx context.Context, prefix ...string) (*Metrics, error) {
	if ctx != nil {
		if metrics, ok := ctx.Value(metricsKey{}).(*Metrics); ok {
			return metrics, nil
		}

		return nil, ErrNoMetrics
	}

	return nil, ErrNilContext
}

// IntoContext returns a new context derived from the provided context which
// carries the provided Metrics.
func IntoContext(ctx context.Context, metrics *Metrics) context.Context {
	return context.WithValue(ctx, metricsKey{}, metrics)
}
