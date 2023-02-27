// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdinreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

// This file implements factory for stdin receiver.

const (
	// The value of "type" key in configuration.
	typeStr        = "stdin"
	stabilityLevel = component.StabilityLevelAlpha
)

type Config struct {
	StdinClosedHook func()
}

// NewFactory creates a factory for stdin receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(newLogsReceiver, stabilityLevel))
}

// CreateDefaultConfig creates the default configuration for stdin receiver.
func createDefaultConfig() component.Config {
	return &Config{}
}
