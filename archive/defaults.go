// Copyright (C) 2021 Talos, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package archive

import (
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	logging "github.com/op/go-logging"

	"time"
)

// FIXME: Provide a method to return the defaults and initialize with parameterised version

type Defaults struct {
	ArchiveLoglevel       logging.Level
	AeronLoglevel         logging.Level
	ControlTimeout        time.Duration      // How long to try sending control messages [see Proxy.Timeout]
	ControlIdleStrategy   idlestrategy.Idler // Idlestrategy as for aeron itself
	RecordingIdleStrategy idlestrategy.Idler // Idlestrategy as for aeron itself
	ControlRetries        int                // How many retries for control messages [see Proxy.Retries] // FIXME: use
	RangeChecking         bool               // Turn on range checking for control protocol marshalling
	ResponseChannel       string             // control response subscription channel
	ResponseStream        int32              // and stream
	RequestChannel        string             // control request publication channel
	RequestStream         int32              // and stream
}

var ArchiveDefaults Defaults = Defaults{
	ArchiveLoglevel:       logging.INFO,
	AeronLoglevel:         logging.INFO,
	ControlTimeout:        time.Second * 5,
	ControlIdleStrategy:   idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
	RecordingIdleStrategy: idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
	ControlRetries:        4,
	RangeChecking:         true, // FIXME: turn off
	ResponseChannel:       "aeron:udp?endpoint=localhost:8020",
	ResponseStream:        20,
	RequestChannel:        "aeron:udp?endpoint=localhost:8010",
	RequestStream:         10,
}
