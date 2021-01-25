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
	"time"
)

type Defaults struct {
	ControlTimeout      time.Duration      // How long to try sending control messages [see Proxy.Timeout]
	ControlRetries      int                // How many retries for control messages [see Proxy.Retries]
	ControlIdleStrategy idlestrategy.Idler // How we wait [various]
	RangeChecking       bool               // Turn on range checking for control protocol marshalling
}

var defaults Defaults = Defaults{
	ControlTimeout:      time.Second * 5,
	ControlRetries:      4,
	ControlIdleStrategy: idlestrategy.Sleeping{SleepFor: time.Millisecond * 5},
	RangeChecking:       true, // FIXME: turn off
}
