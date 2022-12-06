/*
Copyright 2016 Stanislav Liberman
Copyright 2022 Steven Stern

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aeron

import (
	"time"

	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
)

// Context configuration options are located here https://github.com/real-logic/Aeron/wiki/Configuration-Options#aeron-client-options
type Context struct {
	aeronDir      string
	mediaDriverTo time.Duration

	errorHandler func(error)

	newPublicationHandler     NewPublicationHandler
	newSubscriptionHandler    NewSubscriptionHandler
	availableImageHandler     AvailableImageHandler
	unavailableImageHandler   UnavailableImageHandler
	availableCounterHandler   AvailableCounterHandler
	unavailableCounterHandler UnavailableCounterHandler

	resourceLingerTo        time.Duration
	publicationConnectionTo time.Duration
	interServiceTo          time.Duration

	idleStrategy idlestrategy.Idler
}

// NewContext creates and initializes new Context for Aeron
func NewContext() *Context {
	ctx := new(Context)

	ctx.aeronDir = DefaultAeronDir + "/aeron-" + UserName

	ctx.errorHandler = func(err error) { logger.Error(err) }

	ctx.newPublicationHandler = func(string, int32, int32, int64) {}
	ctx.newSubscriptionHandler = func(string, int32, int64) {}
	ctx.availableImageHandler = func(Image) {}
	ctx.unavailableImageHandler = func(Image) {}

	ctx.mediaDriverTo = time.Second * 5
	ctx.resourceLingerTo = time.Second * 3
	ctx.publicationConnectionTo = time.Second * 5
	ctx.interServiceTo = time.Second * 10

	ctx.idleStrategy = idlestrategy.Sleeping{SleepFor: time.Millisecond * 4}

	return ctx
}

// ErrorHandler sets the error handler callback
func (ctx *Context) ErrorHandler(handler func(error)) *Context {
	ctx.errorHandler = handler
	return ctx
}

// AeronDir sets the root directory for media driver files
func (ctx *Context) AeronDir(dir string) *Context {
	ctx.aeronDir = dir
	return ctx
}

// MediaDriverTimeout sets the timeout for keep alives to media driver
func (ctx *Context) MediaDriverTimeout(to time.Duration) *Context {
	ctx.mediaDriverTo = to
	return ctx
}

// ResourceLingerTimeout sets the timeout for resource cleanup after they're released
func (ctx *Context) ResourceLingerTimeout(to time.Duration) *Context {
	ctx.resourceLingerTo = to
	return ctx
}

// InterServiceTimeout sets the timeout for client heartbeat
func (ctx *Context) InterServiceTimeout(to time.Duration) *Context {
	ctx.interServiceTo = to
	return ctx
}

func (ctx *Context) PublicationConnectionTimeout(to time.Duration) *Context {
	ctx.publicationConnectionTo = to
	return ctx
}

// newSubscriptionHandler sets an optional callback for new subscriptions
func (ctx *Context) NewSubscriptionHandler(handler func(string, int32, int64)) *Context {
	ctx.newSubscriptionHandler = handler
	return ctx
}

// newPublicationHandler sets an optional callback for new publications
func (ctx *Context) NewPublicationHandler(handler func(string, int32, int32, int64)) *Context {
	ctx.newPublicationHandler = handler
	return ctx
}

// AvailableImageHandler sets an optional callback for available image notifications
func (ctx *Context) AvailableImageHandler(handler func(Image)) *Context {
	ctx.availableImageHandler = handler
	return ctx
}

// UnavailableImageHandler sets an optional callback for unavailable image notification
func (ctx *Context) UnavailableImageHandler(handler func(Image)) *Context {
	ctx.unavailableImageHandler = handler
	return ctx
}

// CncFileName returns the name of the Counters file
func (ctx *Context) CncFileName() string {
	return ctx.aeronDir + "/" + counters.CncFile
}

// IdleStrategy provides an IdleStrategy for the thread responsible for communicating
// with the Aeron Media Driver.
func (ctx *Context) IdleStrategy(idleStrategy idlestrategy.Idler) *Context {
	ctx.idleStrategy = idleStrategy
	return ctx
}

// AvailableCounterHandler sets up a callback for when a Counter is available.  This will be added to the list before
// additional handlers are added with Aeron.AddAvailableCounterHandler.
func (ctx *Context) AvailableCounterHandler(handler AvailableCounterHandler) {
	ctx.availableCounterHandler = handler
}

// GetAvailableCounterHandler gets the callback handler for when a counter is available.
func (ctx *Context) GetAvailableCounterHandler() AvailableCounterHandler {
	return ctx.availableCounterHandler
}

// UnavailableCounterHandler sets up a callback for when a Counter is unavailable.  This will be added to the list first
// before additional handlers are added with Aeron.AddUnavailableCounterHandler.
func (ctx *Context) UnavailableCounterHandler(handler UnavailableCounterHandler) {
	ctx.unavailableCounterHandler = handler
}

func (ctx *Context) GetUnavailableCounterHandler() UnavailableCounterHandler {
	return ctx.unavailableCounterHandler
}
