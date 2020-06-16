// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package server provides an opinionated http server.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/google/exposure-notifications-server/internal/logging"
)

// ServeMetricsIfPrometheus serves the opencensus metrics at /metrics when OBSERVABILITY_EXPORTER set to "prometheus"
func ServeMetricsIfPrometheus(ctx context.Context, errCh chan error) error {
	logger := logging.FromContext(ctx)

	exporter := os.Getenv("OBSERVABILITY_EXPORTER")
	metricsPort := os.Getenv("METRICS_PORT")
	if strings.EqualFold(exporter, "prometheus") {
		if metricsPort == "" {
			return fmt.Errorf("OBSERVABILITY_EXPORTER set to 'prometheus' but no METRICS_PORT set")
		}
		go func() {
			exporter, err := prometheus.NewExporter(prometheus.Options{})
			if err != nil {
				logger.Debugf("failed to create prometheus exporter: %v", err)
				select {
				case errCh <- err:
				default:
				}
			}

			http.Handle("/metrics", exporter)

			logger.Debugf("listening on :%s", metricsPort)
			if http.ListenAndServe(":"+metricsPort, nil) != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}
	return nil
}
