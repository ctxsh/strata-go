package apex

import "github.com/prometheus/client_golang/prometheus"

type Labels prometheus.Labels

func (l Labels) Keys() []string {
	keys := make([]string, 0)
	for k := range l {
		keys = append(keys, k)
	}
	return keys
}
