// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"github.com/atoulme/stdinotel/receiver/stdinreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/telemetry"
	"go.uber.org/zap/zapcore"
	"os"
)

type configProvider struct {
	collector *otelcol.Collector
}

func (c *configProvider) Get(ctx context.Context, factories otelcol.Factories) (*otelcol.Config, error) {
	cfg, id := createExporterConfig(factories)

	return &otelcol.Config{
		Receivers: map[component.ID]component.Config{
			component.NewID("stdin"): &stdinreceiver.Config{
				StdinClosedHook: func() {
					c.collector.Shutdown()
				},
			},
		},
		Processors: map[component.ID]component.Config{},
		Exporters: map[component.ID]component.Config{
			component.NewID(id): cfg,
		},
		Connectors: map[component.ID]component.Config{},
		Extensions: map[component.ID]component.Config{},
		Service: service.Config{
			Telemetry: telemetry.Config{
				Metrics: telemetry.MetricsConfig{Level: configtelemetry.LevelNone},
				Logs:    telemetry.LogsConfig{Encoding: "console", Level: zapcore.ErrorLevel},
			},
			Extensions: nil,
			Pipelines: map[component.ID]*service.PipelineConfig{
				component.NewID("logs"): {
					Receivers: []component.ID{
						component.NewID("stdin"),
					},
					Processors: []component.ID{},
					Exporters: []component.ID{
						component.NewID(id),
					},
				},
			},
		},
	}, nil
}

func createExporterConfig(factories otelcol.Factories) (component.Config, component.Type) {
	protocol := getProtocol()
	cfg := factories.Exporters[protocol].CreateDefaultConfig()
	if splunkCfg, ok := cfg.(*splunkhecexporter.Config); ok {
		splunkCfg.Endpoint = os.Getenv("STDINOTEL_ENDPOINT")
		splunkCfg.Token = configopaque.String(os.Getenv("STDINOTEL_TOKEN"))
		splunkCfg.Index = os.Getenv("STDINOTEL_SPLUNK_INDEX")
		splunkCfg.TLSSetting.InsecureSkipVerify = os.Getenv("STDINOTEL_TLS_INSECURE_SKIP_VERIFY") == "true"

	}
	if splunkCfg, ok := cfg.(*otlpexporter.Config); ok {
		splunkCfg.Endpoint = os.Getenv("STDINOTEL_ENDPOINT")
	}
	if splunkCfg, ok := cfg.(*otlphttpexporter.Config); ok {
		splunkCfg.Endpoint = os.Getenv("STDINOTEL_ENDPOINT")
	}

	return cfg, protocol
}

func getProtocol() component.Type {
	protocol := os.Getenv("STDINOTEL_PROTOCOL")
	switch protocol {
	case "splunk_hec":
		return "splunk_hec"
	case "otlp":
		return "otlp"
	case "otlphttp":
		return "otlphttp"
	default:
		return "otlp"
	}
}

func (c configProvider) Watch() <-chan error {
	return nil
}

func (c configProvider) Shutdown(_ context.Context) error {
	return nil
}
