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
	"testing"

	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
)

func TestCreateReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	mockLogsConsumer := consumertest.NewNop()
	lReceiver, err := newLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, mockLogsConsumer)
	assert.Nil(t, err, "receiver creation failed")
	assert.NotNil(t, lReceiver, "receiver creation failed")
}
