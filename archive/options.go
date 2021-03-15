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

// These are the Options used by default for an Archive object
// Those attributes marked [init] must be changed before calling ArchiveConnect()
// Those attributes marked [runtime] may be changed at any time
// Those attributes marked [enable] may be changed when the feature is not enabled
type Options struct {
	ArchiveLoglevel        logging.Level      // [runtime]
	AeronLoglevel          logging.Level      // [runtime]
	ControlTimeout         time.Duration      // [runtime] How long to try sending/receiving control messages
	ControlIdleStrategy    idlestrategy.Idler // [runtime] Idlestrategy for sending/receiving control messages
	RangeChecking          bool               // [init] Turn on range checking for control protocol marshalling
	RequestChannel         string             // [init] Control request publication channel
	RequestStream          int32              // [init] and stream
	ResponseChannel        string             // [init] Control response subscription channel
	ResponseStream         int32              // [init] and stream
	RecordingEventsChannel string             // [enable] Recording progress events
	RecordingEventsStream  int32              // [enable] and stream
	RecordingIdleStrategy  idlestrategy.Idler // [runtime] Idlestrategy for the recording poller
}

var ArchiveOptions Options = Options{
	ArchiveLoglevel:        logging.INFO,
	AeronLoglevel:          logging.NOTICE,
	ControlTimeout:         time.Second * 5,
	ControlIdleStrategy:    idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
	RangeChecking:          true,                                                   // FIXME: turn off
	RequestChannel:         "aeron:udp?endpoint=localhost:8010",
	RequestStream:          10,
	ResponseChannel:        "aeron:udp?endpoint=localhost:8020",
	ResponseStream:         20,
	RecordingEventsChannel: "aeron:udp?control-mode=dynamic|control=localhost:8030",
	RecordingEventsStream:  30,
	RecordingIdleStrategy:  idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
}
