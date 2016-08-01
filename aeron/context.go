/*
Copyright 2016 Stanislav Liberman

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
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"os/user"
	"time"
)

// https://github.com/real-logic/Aeron/wiki/Configuration-Options#aeron-client-options
type Context struct {
	aeronDir      string // aeron.dir
	mediaDriverTo time.Duration

	errorHandler func(error)

	newPublicationHandler   NewPublicationHandler
	newSubscriptionHandler  NewSubscriptionHandler
	availableImageHandler   AvailableImageHandler
	unavailableImageHandler UnavailableImageHandler

	resourceLingerTo        time.Duration
	publicationConnectionTo time.Duration
	interServiceTo          time.Duration

	idleStrategy idlestrategy.Idler
}

func NewContext() *Context {
	ctx := new(Context)

	ctx.aeronDir = "/tmp"

	ctx.errorHandler = func(err error) { logger.Error(err) }

	ctx.newPublicationHandler = func(string, int32, int32, int64) {}
	ctx.newSubscriptionHandler = func(string, int32, int64) {}
	ctx.availableImageHandler = func(*Image) {}
	ctx.unavailableImageHandler = func(*Image) {}

	ctx.mediaDriverTo = time.Second
	ctx.resourceLingerTo = time.Second * 3
	ctx.publicationConnectionTo = time.Second * 5
	ctx.interServiceTo = time.Second * 10

	ctx.idleStrategy = idlestrategy.Sleeping{SleepFor: time.Millisecond * 4}

	return ctx
}

func (ctx *Context) ErrorHandler(handler func(error)) *Context {
	ctx.errorHandler = handler
	return ctx
}

func (ctx *Context) AeronDir(dir string) *Context {
	ctx.aeronDir = dir
	return ctx
}

func (ctx *Context) MediaDriverTimeout(to time.Duration) *Context {
	ctx.mediaDriverTo = to
	return ctx
}

func (ctx *Context) ResourceLingerTimeout(to time.Duration) *Context {
	ctx.resourceLingerTo = to
	return ctx
}

func (ctx *Context) InterServiceTimeout(to time.Duration) *Context {
	ctx.interServiceTo = to
	return ctx
}

func (ctx *Context) PublicationConnectionTimeout(to time.Duration) *Context {
	ctx.publicationConnectionTo = to
	return ctx
}

func (ctx *Context) AvailableImageHandler(handler func(*Image)) *Context {
	ctx.availableImageHandler = handler
	return ctx
}

func (ctx *Context) UnavailableImageHandler(handler func(*Image)) *Context {
	ctx.unavailableImageHandler = handler
	return ctx
}

func (ctx *Context) cncFileName() string {
	user, err := user.Current()
	var uName string = "unknown"
	if err != nil {
		logger.Warningf("Failed to get current user name: %v", err)
	}
	uName = user.Username
	return ctx.aeronDir + "/aeron-" + uName + "/" + counters.Descriptor.CNC_FILE
}
