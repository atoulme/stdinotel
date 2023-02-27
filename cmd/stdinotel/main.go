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

// Program stdinotel is an OpenTelemetry Collector binary.
// This code was first generated, then modified for the purpose of registering a shutdown hook on the stdin receiver.
package main

import (
	"github.com/spf13/cobra"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"
	"log"
)

func main() {
	factories, err := components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "stdinotel",
		Description: "Standard input to OTLP/HEC collector",
		Version:     "0.1.0-dev",
	}

	set := otelcol.CollectorSettings{BuildInfo: info, Factories: factories, ConfigProvider: &configProvider{}}

	if err := run(set); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(params otelcol.CollectorSettings) error {
	cmd := createCommand(params)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("collector server run finished with error: %v", err)
	}

	return nil
}

func createCommand(set otelcol.CollectorSettings) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          set.BuildInfo.Command,
		Version:      set.BuildInfo.Version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			col, err := otelcol.NewCollector(set)
			if err != nil {
				return err
			}
			set.ConfigProvider.(*configProvider).collector = col
			return col.Run(cmd.Context())
		},
	}
	return rootCmd
}
