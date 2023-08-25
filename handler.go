package strata

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HandlerFor returns the handler for the metrics registry.
func HandlerFor(metrics *Metrics) http.Handler {
	return promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{
		Timeout: DefaultTimeout,
	})
}
