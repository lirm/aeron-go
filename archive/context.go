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
	"time"
)

// The ArchiveContext contains data (Listeners, Options, aeronContext, etc) which is
// necessary across the archive, in use in the Proxy, Control, and RecordingEvent.
//
// It also wraps some the Aeron context methods so you can deal with just the archive
type ArchiveContext struct {
	Options      *Options
	SessionId    int64
	aeron        *aeron.Aeron
	aeronContext *aeron.Context
}

// NewContext creates and initializes a new Context for AeronArchive
// FIXME: For now this is just a simple wrapperarc, round this out
// FIXME: ArchiveContext configuration options are located here:
// https://github.com/real-logic/Aeron/wiki/Configuration-Options#aeron-archive-client-opt
func NewArchiveContext() *ArchiveContext {

	context := new(ArchiveContext)
	context.aeronContext = aeron.NewContext()

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

// CncFileName returns the name of the Counters file
func (context *ArchiveContext) CncFileName() string {
	return context.aeronContext.CncFileName()
}
