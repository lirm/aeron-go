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
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"time"
)

// FIXME: ArchiveContext configuration options are located here:
// https://github.com/real-logic/Aeron/wiki/Configuration-Options#aeron-archive-client-options
type ArchiveContext struct {
	aeron                   *aeron.Aeron
	aeronContext            *aeron.Context
	RecordingEventsChannel  string
	RecordingEventsStreamID int32
	MessageTimeouts         time.Duration
	IdleStrategy            idlestrategy.Idler
}

// NewContext creates and initializes new Context for AeronArchive
// FIXME: For now this is just a simple wrapper, round this out
// FIXME: ArchiveContext configuration options are located here:
// https://github.com/real-logic/Aeron/wiki/Configuration-Options#aeron-archive-client-opt
func NewArchiveContext() *ArchiveContext {

	context := new(ArchiveContext)
	context.aeronContext = aeron.NewContext()

	// FIXME: Only in debug
	context.aeronContext.AvailableImageHandler(ArchiveAvailableImageHandler)
	context.aeronContext.UnavailableImageHandler(ArchiveUnavailableImageHandler)
	context.aeronContext.NewSubscriptionHandler(ArchiveNewSubscriptionHandler)
	context.aeronContext.NewPublicationHandler(ArchiveProxyNewPublicationHandler)

	// Archive specific additional context
	// FIXME: Add methods to set all of these and make a suitable place for the defaults
	// FIXME: Java/C++ try and work out a response channel

	context.MessageTimeouts = 5 * 1000 * 1000 * 1000                             // 5 seconds in Nanos
	context.IdleStrategy = idlestrategy.Sleeping{SleepFor: time.Millisecond * 5} // FIXME: load from defaults

	return context
}

// ErrorHandler sets the error handler callback
func (context *ArchiveContext) ErrorHandler(handler func(error)) *ArchiveContext {
	context.aeronContext.ErrorHandler(handler)
	return context
}

// AeronDir sets the root directory for media driver files
func (context *ArchiveContext) AeronDir(dir string) *ArchiveContext {
	context.aeronContext.AeronDir(dir)
	return context
}

// MediaDriverTimeout sets the timeout for keep alives to media driver
func (context *ArchiveContext) MediaDriverTimeout(timeout time.Duration) *ArchiveContext {
	context.aeronContext.MediaDriverTimeout(timeout)
	return context
}

// ResourceLingerTimeout sets the timeout for resource cleanup after they're released
func (context *ArchiveContext) ResourceLingerTimeout(timeout time.Duration) *ArchiveContext {
	context.aeronContext.ResourceLingerTimeout(timeout)
	return context
}

// InterServiceTimeout sets the timeout for client heartbeat
func (context *ArchiveContext) InterServiceTimeout(timeout time.Duration) *ArchiveContext {
	context.aeronContext.InterServiceTimeout(timeout)
	return context
}

// PublicationConnectionTimeout sets the timeout for publications
func (context *ArchiveContext) PublicationConnectionTimeout(timeout time.Duration) *ArchiveContext {
	context.aeronContext.PublicationConnectionTimeout(timeout)
	return context
}

// AvailableImageHandler sets an optional callback for available image notifications
func (context *ArchiveContext) AvailableImageHandler(handler func(*aeron.Image)) *ArchiveContext {
	context.aeronContext.AvailableImageHandler(handler)
	return context
}

// UnavailableImageHandler sets an optional callback for unavailable image notification
func (context *ArchiveContext) UnavailableImageHandler(handler func(*aeron.Image)) *ArchiveContext {
	context.aeronContext.UnavailableImageHandler(handler)
	return context
}

// CncFileName returns the name of the Counters file
func (context *ArchiveContext) CncFileName() string {
	return context.aeronContext.CncFileName()
}
