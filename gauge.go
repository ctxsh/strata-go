package apex

func (m *Metrics) Set(value float64, labels Labels) {}

func (m *Metrics) Add(value float64, labels Labels) {}

func (m *Metrics) Sub(value float64, labels Labels) {}

// func (m *Metrics) getGauge(name string) {}

// func (m *Metrics) registerGauge(name string, labels []string) (*prometheus.GaugeVec, error) {}
