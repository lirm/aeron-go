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

// Package archive provides API access to Aeron's archive-media-driver
package archive

import (
	"github.com/lirm/aeron-go/aeron/atomic"
)

// Channel to receive a new unique 64 bit int used for SessionIDs and CorrelationIDs
var ID64 chan int64

// Send to it To stop the generator
var ID64Terminate bool

// Start the generator as a goroutine
func init() {
	ID64Terminate = false
	StartIdGenerator()
}

func id64() {
	logger.Debug("Starting ID64 generator")
	var i64 atomic.Long

	for {
		ID64 <- i64.Inc()
		if ID64Terminate {
			return
		}
	}
}

// Start the correlationID generator
func StartIdGenerator() {
	ID64 = make(chan int64, 1)
	go id64()
}

// Stop the correlationID generator
func StopIdGenerator() {
	ID64Terminate = true
	_ = <-ID64
}
