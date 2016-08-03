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
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/ringbuffer"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/driver"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"time"
	"github.com/lirm/aeron-go/aeron/atomic"
)

type NewPublicationHandler func(string, int32, int32, int64)

type NewSubscriptionHandler func(string, int32, int64)

type AvailableImageHandler func(*Image)

type UnavailableImageHandler func(*Image)

type Aeron struct {
	context            *Context
	conductor          ClientConductor
	toDriverRingBuffer rb.ManyToOne
	driverProxy        driver.Proxy

	toDriverAtomicBuffer  *atomic.Buffer
	toClientsAtomicBuffer *atomic.Buffer
	counterValuesBuffer   *atomic.Buffer

	cncBuffer *memmap.File

	toClientsBroadcastReceiver *broadcast.Receiver
	toClientsCopyReceiver      *broadcast.CopyReceiver
}

var logger = logging.MustGetLogger("aeron")

func Connect(ctx *Context) *Aeron {
	aeron := new(Aeron)
	aeron.context = ctx
	logger.Debugf("Connecting with context: %v", ctx)

	aeron.cncBuffer = counters.MapFile(ctx.cncFileName())

	aeron.toDriverAtomicBuffer = counters.CreateToDriverBuffer(aeron.cncBuffer)
	aeron.toClientsAtomicBuffer = counters.CreateToClientsBuffer(aeron.cncBuffer)
	aeron.counterValuesBuffer = counters.CreateCounterValuesBuffer(aeron.cncBuffer)

	aeron.toDriverRingBuffer.Init(aeron.toDriverAtomicBuffer)

	aeron.driverProxy.Init(&aeron.toDriverRingBuffer)

	aeron.toClientsBroadcastReceiver = broadcast.NewReceiver(aeron.toClientsAtomicBuffer)

	aeron.toClientsCopyReceiver = broadcast.NewCopyReceiver(aeron.toClientsBroadcastReceiver)

	clientLivenessTo := time.Duration(counters.ClientLivenessTimeout(aeron.cncBuffer))

	aeron.conductor.Init(&aeron.driverProxy, aeron.toClientsCopyReceiver, clientLivenessTo, ctx.mediaDriverTo,
		ctx.publicationConnectionTo, ctx.resourceLingerTo)
	aeron.conductor.counterValuesBuffer = aeron.counterValuesBuffer

	aeron.conductor.onAvailableImageHandler = ctx.availableImageHandler
	aeron.conductor.onUnavailableImageHandler = ctx.unavailableImageHandler

	go aeron.conductor.Run(ctx.idleStrategy)

	return aeron
}

func (aeron *Aeron) Close() error {
	err := aeron.conductor.Close()

	err = aeron.cncBuffer.Close()

	return err
}

func (aeron *Aeron) AddSubscription(channel string, streamId int32) chan *Subscription {
	ch := make(chan *Subscription, 1)

	regId := aeron.conductor.AddSubscription(channel, streamId)
	go func() {
		subscription := aeron.conductor.FindSubscription(regId)
		for subscription == nil {
			subscription = aeron.conductor.FindSubscription(regId)
			if subscription == nil {
				aeron.context.idleStrategy.Idle(0)
			}
		}
		ch <- subscription
		close(ch)
	}()

	return ch
}

func (aeron *Aeron) AddPublication(channel string, streamId int32) chan *Publication {
	ch := make(chan *Publication, 1)

	regId := aeron.conductor.AddPublication(channel, streamId)
	go func() {
		publication := aeron.conductor.FindPublication(regId)
		for publication == nil {
			publication = aeron.conductor.FindPublication(regId)
			if publication == nil {
				aeron.context.idleStrategy.Idle(0)
			}
		}
		ch <- publication
		close(ch)
	}()

	return ch
}
