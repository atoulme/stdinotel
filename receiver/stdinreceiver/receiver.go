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
	"bufio"
	"context"
	"errors"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"os"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

const (
	transport = "stdin"
	format    = "string"
)

var (
	errNilNextLogsConsumer = errors.New("nil logsConsumer")
	listenerEnabled        = sync.Once{}
	listeners              []chan string
	stdin                  = os.Stdin
)

// stdinReceiver implements the component.MetricsReceiver for stdin metric protocol.
type stdinReceiver struct {
	logger       *zap.Logger
	config       *Config
	logsConsumer consumer.Logs
	shutdown     chan struct{}
}

// NewLogsReceiver creates the stdin receiver with the given configuration.
func NewLogsReceiver(
	logger *zap.Logger,
	config Config,
	nextConsumer consumer.Logs,
) (receiver.Logs, error) {
	if nextConsumer == nil {
		return nil, errNilNextLogsConsumer
	}

	r := &stdinReceiver{
		logger:       logger,
		config:       &config,
		logsConsumer: nextConsumer,
	}

	return r, nil
}

func startStdinListener(logger *zap.Logger) {
	listenerEnabled.Do(func() {
		reader := bufio.NewReader(stdin)
		for {
			scanner := bufio.NewScanner(reader)
			scanner.Split(bufio.ScanLines) // Set up the split function.
			for scanner.Scan() {
				line := scanner.Text()
				for _, listener := range listeners {
					listener <- line
				}
			}
			if err := scanner.Err(); err != nil {
				logger.Error("Error reading stdin", zap.Error(err))
			}
		}
	})
}

// Start starts the stdin receiver, adding it to the stdin listener.
func (r *stdinReceiver) Start(ctx context.Context, host component.Host) error {
	r.shutdown = make(chan struct{})

	listener := make(chan string)
	listeners = append(listeners, listener)

	go func() {
		go startStdinListener(r.logger)
		for {
			select {
			case nextLine := <-listener:
				err := r.consumeLine(ctx, nextLine)
				if err != nil {
					r.logger.Error("error with log", zap.Error(err))
				}
			case <-r.shutdown:
				return
			}
		}
	}()

	return nil
}

func (r *stdinReceiver) consumeLine(ctx context.Context, line string) error {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	sl.LogRecords().AppendEmpty().Body().SetStr(line)
	err := r.logsConsumer.ConsumeLogs(ctx, ld)
	return err
}

// Shutdown shuts down the stdin receiver, closing its listener to the stdin loop.
func (r *stdinReceiver) Shutdown(context.Context) error {
	close(r.shutdown)

	return nil
}
