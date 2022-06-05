package utils

import (
	"strings"

	"github.com/ctxswitch/apex-go/errors"
	"github.com/prometheus/client_golang/prometheus"
)

func NameBuilder(ns string, sub string, name string, sep rune) (string, error) {
	var builder strings.Builder

	if ns != "" {
		builder.WriteString(ns)
		builder.WriteRune(sep)
	}

	if sub != "" {
		builder.WriteString(sub)
		builder.WriteRune(sep)
	}

	if name == "" {
		return "", errors.ErrInvalidMetricName
	}
	builder.WriteString(name)
	return builder.String(), nil
}

func LabelKeys(labels prometheus.Labels) []string {
	keys := make([]string, 0)
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}

func Register(metric prometheus.Collector) error {
	if err := prometheus.Register(metric); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return errors.ErrAlreadyRegistered
		} else {
			return errors.ErrRegistrationFailed
		}
	}
	return nil
}
