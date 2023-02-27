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
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/multierr"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var (
	stdin = os.Stdin
)

// stdinReceiver implements the component.MetricsReceiver for stdin metric protocol.
type stdinReceiver struct {
	logger       *zap.Logger
	config       *Config
	logsConsumer consumer.Logs
	obsrecv      *obsreport.Receiver
	wg           sync.WaitGroup
}

// newLogsReceiver creates the stdin receiver with the given configuration.
func newLogsReceiver(
	ctx context.Context,
	settings receiver.CreateSettings,
	config component.Config,
	nextConsumer consumer.Logs,
) (receiver.Logs, error) {

	cfg, ok := config.(*Config)
	if !ok {
		return nil, errors.New("invalid config")
	}

	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             settings.ID,
		Transport:              "",
		ReceiverCreateSettings: settings,
	})
	if err != nil {
		return nil, err
	}

	r := &stdinReceiver{
		logger:       settings.Logger,
		config:       cfg,
		logsConsumer: nextConsumer,
		obsrecv:      obsrecv,
	}

	return r, nil
}

func (r *stdinReceiver) startStdinListener(ctx context.Context, logger *zap.Logger, host component.Host) {
	r.obsrecv.StartLogsOp(ctx)
	var errs []error
	i := 0
	reader := bufio.NewReader(stdin)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines) // Set up the split function.
	for scanner.Scan() {
		line := scanner.Text()
		err := r.consumeLine(ctx, line)
		if err != nil {
			errs = append(errs, err)
		} else {
			i++
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error reading stdin", zap.Error(err))
	}
	combined := multierr.Combine(errs...)
	r.obsrecv.EndLogsOp(ctx, "", i, combined)
	if len(errs) != 0 {
		host.ReportFatalError(combined)
	} else {
		r.config.StdinClosedHook()
	}
	r.wg.Done()
}

// Start starts the stdin receiver.
func (r *stdinReceiver) Start(ctx context.Context, host component.Host) error {
	r.wg.Add(1)
	go r.startStdinListener(ctx, r.logger, host)
	return nil
}

func (r *stdinReceiver) consumeLine(ctx context.Context, line string) error {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	lr := sl.LogRecords().AppendEmpty()
	lr.Body().SetStr(line)
	lr.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	err := r.logsConsumer.ConsumeLogs(ctx, ld)
	return err
}

// Shutdown shuts down the stdin receiver.
func (r *stdinReceiver) Shutdown(context.Context) error {
	r.wg.Wait()
	return nil
}
