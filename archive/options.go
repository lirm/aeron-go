// Copyright (C) 2021-2022 Talos, Inc.
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
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"go.uber.org/zap/zapcore"
	"time"
)

// Options are set by NewArchive() with the default options specified below
//
// Some options may be changed dynamically by setting their values in Archive.Context.Options.*
// Those attributes marked [compile] must be changed at compile time
// Those attributes marked [init] must be changed before calling ArchiveConnect()
// Those attributes marked [runtime] may be changed at any time
// Those attributes marked [enable] may be changed when the feature is not enabled
type Options struct {
	RequestChannel         string             // [init] Control request publication channel
	RequestStream          int32              // [init] and stream
	ResponseChannel        string             // [init] Control response subscription channel
	ResponseStream         int32              // [init] and stream
	RecordingEventsChannel string             // [enable] Recording progress events
	RecordingEventsStream  int32              // [enable] and stream
	ArchiveLoglevel        zapcore.Level      // [runtime] via logging.SetLevel()
	AeronLoglevel          zapcore.Level      // [runtime] via logging.SetLevel()
	Timeout                time.Duration      // [runtime] How long to try sending/receiving control messages
	IdleStrategy           idlestrategy.Idler // [runtime] Idlestrategy for sending/receiving control messages
	RangeChecking          bool               // [runtime] archive protocol marshalling checks
	AuthEnabled            bool               // [init] enable to require AuthConnect() over Connect()
	AuthCredentials        []uint8            // [init] The credentials to be provided to AuthConnect()
	AuthChallenge          []uint8            // [init] The challenge string we are to expect (checked iff not nil)
	AuthResponse           []uint8            // [init] The challengeResponse we should provide
}

// These are the Options used by default for an Archive object
var defaultOptions = Options{
	RequestChannel:         "aeron:udp?endpoint=localhost:8010",
	RequestStream:          10,
	ResponseChannel:        "aeron:udp?endpoint=localhost:8020",
	ResponseStream:         20,
	RecordingEventsChannel: "aeron:udp?control-mode=dynamic|control=localhost:8030",
	RecordingEventsStream:  30,
	ArchiveLoglevel:        zapcore.WarnLevel,
	AeronLoglevel:          zapcore.WarnLevel,
	Timeout:                time.Second * 5,
	IdleStrategy:           idlestrategy.Sleeping{SleepFor: time.Millisecond * 50}, // FIXME: tune
	RangeChecking:          false,
	AuthEnabled:            false,
	AuthCredentials:        nil,
	AuthChallenge:          nil,
	AuthResponse:           nil,
}

// DefaultOptions creates and returns a new Options from the defaults.
func DefaultOptions() *Options {
	options := defaultOptions
	return &options
}
