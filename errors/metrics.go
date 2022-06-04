package errors

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type ApexInternalErrorMetrics struct {
	errPanicRecovery      *prometheus.CounterVec
	errInvalidMetricName  *prometheus.CounterVec
	errRegistrationFailed *prometheus.CounterVec
	errAlreadyRegistered  *prometheus.CounterVec
}

func NewApexInternalErrorMetrics(ns string, sub string, sep rune) *ApexInternalErrorMetrics {
	var builder strings.Builder

	if ns != "" {
		builder.WriteString(ns)
		builder.WriteRune(sep)
	}

	if sub != "" {
		builder.WriteString(sub)
		builder.WriteRune(sep)
	}

	builder.WriteString("apex")
	builder.WriteRune(sep)
	builder.WriteString("error")
	builder.WriteRune(sep)

	prefix := builder.String()

	errPanicRecovery := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "panic_recovery",
	}, []string{"name", "type"})

	errInvalidMetricName := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "invalid_metric_name",
	}, []string{"name", "type"})

	errRegistrationFailed := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "registration_failed",
	}, []string{"name", "type"})

	errAlreadyRegistered := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "already_registered",
	}, []string{"name", "type"})

	_ = register(errPanicRecovery)
	_ = register(errInvalidMetricName)
	_ = register(errRegistrationFailed)
	_ = register(errAlreadyRegistered)

	return &ApexInternalErrorMetrics{
		errPanicRecovery:      errPanicRecovery,
		errInvalidMetricName:  errInvalidMetricName,
		errRegistrationFailed: errRegistrationFailed,
		errAlreadyRegistered:  errAlreadyRegistered,
	}
}

func (a *ApexInternalErrorMetrics) PanicRecovery(name string, t string) {
	a.errPanicRecovery.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func (a *ApexInternalErrorMetrics) InvalidMetricName(name string, t string) {
	a.errInvalidMetricName.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func (a *ApexInternalErrorMetrics) RegistrationFailed(name string, t string) {
	a.errRegistrationFailed.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func (a *ApexInternalErrorMetrics) AlreadyRegistered(name string, t string) {
	a.errAlreadyRegistered.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func register(metric prometheus.Collector) error {
	if err := prometheus.Register(metric); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return err
		} else {
			panic(err)
		}
	}
	return nil
}
