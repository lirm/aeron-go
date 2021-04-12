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

// Options are set by NewArchive() with default ioptions specified below
// Some options may be changed dynamically by setting their values in Archive.Context.Options.*
// Those attributes marked [compile] must be changed at compile time
// Those attributes marked [init] must be changed before calling ArchiveConnect()
// Those attributes marked [runtime] may be changed at any time
// Those attributes marked [enable] may be changed when the feature is not enabled
type Options struct {
	RangeChecking          bool               // [runtime] archive protocol marshalling checks
	ArchiveLoglevel        logging.Level      // [runtime]
	AeronLoglevel          logging.Level      // [runtime]
	Timeout                time.Duration      // [runtime] How long to try sending/receiving control messages
	IdleStrategy           idlestrategy.Idler // [runtime] Idlestrategy for sending/receiving control messages
	RequestChannel         string             // [init] Control request publication channel
	RequestStream          int32              // [init] and stream
	ResponseChannel        string             // [init] Control response subscription channel
	ResponseStream         int32              // [init] and stream
	RecordingEventsChannel string             // [enable] Recording progress events
	RecordingEventsStream  int32              // [enable] and stream
}

// These are the Options used by default for an Archive object
var defaultOptions Options = Options{
	RangeChecking:          false,
	ArchiveLoglevel:        logging.INFO,
	AeronLoglevel:          logging.NOTICE,
	Timeout:                time.Second * 5,
	IdleStrategy:           idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
	RequestChannel:         "aeron:udp?endpoint=localhost:8010",
	RequestStream:          10,
	ResponseChannel:        "aeron:udp?endpoint=localhost:8020",
	ResponseStream:         20,
	RecordingEventsChannel: "aeron:udp?control-mode=dynamic|control=localhost:8030",
	RecordingEventsStream:  30,
}

// Create and return a new options from the defaults.
func DefaultOptions() *Options {
	options := defaultOptions
	return &options
}
