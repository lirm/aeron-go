// Copyright 2016 Stanislav Liberman
// Copyright 2022 Steven Stern
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

package aeron

import "github.com/lirm/aeron-go/aeron/logbuffer/term"

// Image is a Java-style interface for the image struct.  This is to allow dependency injection and testing of
// the many structs that use image, without deviating from the existing function signatures and code structure.
type Image interface {
	IsClosed() bool
	Poll(handler term.FragmentHandler, fragmentLimit int) int
	BoundedPoll(handler term.FragmentHandler, limitPosition int64, fragmentLimit int) int
	ControlledPoll(handler term.ControlledFragmentHandler, fragmentLimit int) int
	Position() int64
	IsEndOfStream() bool
	SessionID() int32
	CorrelationID() int64
	SubscriptionRegistrationID() int64
	TermBufferLength() int32
	ActiveTransportCount() int32
	Close() error
}
