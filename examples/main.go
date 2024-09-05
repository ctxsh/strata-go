// Copyright (C) 2022, Rob Lyon <rob@ctxswitch.com>
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package main

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"ctx.sh/strata"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func random(min int, max int) float64 {
	return float64(min) + rand.Float64()*(float64(max-min))
}

func runOnce(ctx context.Context) error {
	m, err := strata.FromContext(ctx)
	if err != nil {
		return err
	}
	// Histogram timer
	timer := m.HistogramTimer("latency")
	defer timer.ObserveDuration()

	// Counter functions
	n := m.WithPrefix("func").WithLabels("label", "another")
	n.CounterInc("test_counter", "value1", "another1")
	n.CounterAdd("test_counter", 5.0, "value1", "another2")
	n.CounterInc("test_counter", "value2", "another1")
	n.CounterAdd("test_counter", 2.0, "value2", "another2")

	// Gauge functions
	n.GaugeInc("test_gauge", "value2", "another1")
	n.GaugeSet("test_gauge", random(1, 100), "value2", "another1")
	n.GaugeAdd("test_gauge", 2.0, "value2", "another1")
	n.GaugeSub("test_gauge", 1.0, "value2", "another1")

	// Summary observation
	n.SummaryObserve("test_summary", random(0, 10), "value3", "another1")

	delay := time.Duration(random(1, 1500)) * time.Millisecond
	time.Sleep(delay)

	return nil
}

func main() {
	certFile := flag.String("cert", "", "path to the ssl cert")
	keyFile := flag.String("key", "", "path to the ssl key")
	flag.Parse()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	zapCfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   false,
		Sampling:      nil,
		Encoding:      "console",
		EncoderConfig: encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	zl := zap.Must(zapCfg.Build())
	defer zl.Sync()
	logger := zapr.NewLogger(zl)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	metrics := strata.New(strata.MetricsOpts{
		Logger:         logger,
		Separator:      ':',
		PanicOnError:   true,
		ConstantLabels: []string{"role", "server"},
		SummaryOpts: &strata.SummaryOpts{
			MaxAge:     10 * time.Minute,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			AgeBuckets: 5,
		},
		HistogramBuckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
	}).WithPrefix("strata", "example")

	var obs sync.WaitGroup
	obs.Add(1)
	go func() {
		defer obs.Done()
		logger.Info("starting metrics")
		err := metrics.Start(ctx, strata.ServerOpts{
			Port:                   9090,
			TerminationGracePeriod: 10 * time.Second,
			TLS: &strata.TLSOpts{
				CertFile: *certFile,
				KeyFile:  *keyFile,
			},
		})
		if err != nil {
			panic("could not start metrics")
		}
	}()

	var app sync.WaitGroup
	app.Add(1)
	go func() {
		defer app.Done()
		logger.Info("Starting services")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				metrics.CounterInc("loop_total")
				// This looks odd, but it's not meant to be used this way in normal conditions.
				// I just want to test out the context functions.
				_ = runOnce(strata.IntoContext(ctx, metrics.WithPrefix("runOnce")))
			}
		}
	}()

	<-ctx.Done()
	logger.Info("signal caught, waiting for app to shut down.")
	app.Wait()

	logger.Info("app has shut down, waiting for metrics to shut down.")
	obs.Wait()

	logger.Info("finished")
}
