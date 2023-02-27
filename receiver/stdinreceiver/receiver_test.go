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
	"context"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
)

func TestConsumeLine(t *testing.T) {
	sink := new(consumertest.LogsSink)
	config := createDefaultConfig()
	r := stdinReceiver{logsConsumer: sink, config: config.(*Config)}
	r.consumeLine(context.Background(), "foo")
	lds := sink.AllLogs()
	assert.Equal(t, 1, len(lds))
	log := lds[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Str()
	assert.Equal(t, "foo", log)
}

func TestReceiveLines(t *testing.T) {
	read, write, err := os.Pipe()
	assert.NoError(t, err)
	stdin = read
	sink := new(consumertest.LogsSink)
	config := createDefaultConfig()
	settings := receivertest.NewNopCreateSettings()
	receiverSettings := obsreport.ReceiverSettings{
		ReceiverID:             settings.ID,
		Transport:              "",
		ReceiverCreateSettings: settings,
	}
	obsrecv, _ := obsreport.NewReceiver(receiverSettings)
	r := stdinReceiver{logsConsumer: sink, config: config.(*Config), obsrecv: obsrecv}
	err = r.Start(context.Background(), componenttest.NewNopHost())
	assert.NoError(t, err)
	write.WriteString("foo\nbar\nfoobar\n")
	write.WriteString("foo\r\nbar\nfoobar\n")
	time.Sleep(time.Second * 1)
	lds := sink.AllLogs()
	assert.Equal(t, 6, len(lds))

	read.Chmod(000)
	write.WriteString("foo\nbar\nfoobar\n")
	time.Sleep(time.Second * 1)
	write.Close() // close stdin early.

	err = r.Shutdown(context.Background())
	assert.NoError(t, err)
}
