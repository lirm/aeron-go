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

package driver

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/command"
	rb "github.com/lirm/aeron-go/aeron/ringbuffer"
)

// Proxy is a media driver proxy class that is used to send commands
type Proxy struct {
	toDriverCommandBuffer *rb.ManyToOne
	clientID              int64
}

// Init initializes media driver proxy class
func (driver *Proxy) Init(buffer *rb.ManyToOne) *Proxy {
	driver.toDriverCommandBuffer = buffer
	driver.clientID = driver.toDriverCommandBuffer.NextCorrelationID()

	return driver
}

// TimeOfLastDriverKeepalive gets the time of the last keep alive update sent to media driver
func (driver *Proxy) TimeOfLastDriverKeepalive() int64 {
	return driver.toDriverCommandBuffer.ConsumerHeartbeatTime()
}

// AddSubscription sends driver command to add new subscription
func (driver *Proxy) AddSubscription(channel string, streamID int32) int64 {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddSubscription: correlationId=%d", correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.SubscriptionMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.RegistrationCorrelationID.Set(-1)
		message.StreamID.Set(streamID)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.AddSubscription
	}

	driver.writeCommandToDriver(filler)

	return correlationID

}

// RemoveSubscription sends driver command to remove subscription
func (driver *Proxy) RemoveSubscription(registrationID int64) {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemoveSubscription: correlationId=%d (subId=%d)", correlationID, registrationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.RegistrationID.Set(registrationID)

		*length = message.Size()

		return command.RemoveSubscription
	}

	driver.writeCommandToDriver(filler)
}

// AddPublication sends driver command to add new publication
func (driver *Proxy) AddPublication(channel string, streamID int32) int64 {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddPublication: correlationId=%d", correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.PublicationMessage
		message.Wrap(buffer, 0)
		// For some reason, setting ClientID causes publications to hang.
		// I haven't gotten a chance to look at why.
		// message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.StreamID.Set(streamID)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.AddPublication
	}

	driver.writeCommandToDriver(filler)

	return correlationID
}

// RemovePublication sends driver command to remove publication
func (driver *Proxy) RemovePublication(registrationID int64) {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemovePublication: correlationId=%d (pudId=%d)", correlationID, registrationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.RegistrationID.Set(registrationID)

		*length = message.Size()

		return command.RemovePublication
	}

	driver.writeCommandToDriver(filler)
}

// SendClientKeepalive send keep alive message to the driver
func (driver *Proxy) SendClientKeepalive() {

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.CorrelatedMessage
		message.Wrap(buffer, 0)
		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(0)

		*length = message.Size()

		return command.ClientKeepalive
	}

	driver.writeCommandToDriver(filler)
}

func (driver *Proxy) writeCommandToDriver(filler func(*atomic.Buffer, *int) int32) {
	messageBuffer := make([]byte, 512)

	buffer := atomic.MakeBuffer(messageBuffer)

	length := len(messageBuffer)

	msgTypeID := filler(buffer, &length)

	if !driver.toDriverCommandBuffer.Write(int32(msgTypeID), buffer, 0, int32(length)) {
		panic("couldn't write command to driver")
	}
}
