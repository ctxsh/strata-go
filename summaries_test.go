package apex

// func TestSummarySetMismatchedLabels(t *testing.T) {
// 	m := New(MetricsOpts{MustRegister: true})
// 	m.NewSummary("test_metric_observe", []string{"this", "that"})

// 	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
// 	assert.NotPanics(t, assert.PanicTestFunc(func() {
// 		m.SummaryObserve("test_metric_observe", 100, Labels{"this": "one"})
// 	}))

// 	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
// }

// func TestSummaryAddInvalidCounter(t *testing.T) {
// 	m := New(MetricsOpts{MustRegister: true})
// 	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
// 	assert.NotPanics(t, assert.PanicTestFunc(func() {
// 		m.GaugeAdd("were_is_waldo", 100, Labels{"this": "one"})
// 	}))

// 	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
// }
