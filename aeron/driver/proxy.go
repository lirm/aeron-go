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
	"github.com/lirm/aeron-go/aeron/buffer"
	"github.com/lirm/aeron-go/aeron/buffer/rb"
	"github.com/lirm/aeron-go/aeron/command"
)

type Proxy struct {
	toDriverCommandBuffer *rb.ManyToOne
	clientId              int64
}

func (driver *Proxy) Init(buffer *rb.ManyToOne) *Proxy {
	driver.toDriverCommandBuffer = buffer
	driver.clientId = driver.toDriverCommandBuffer.NextCorrelationId()

	return driver
}

func (driver *Proxy) TimeOfLastDriverKeepalive() int64 {
	return driver.toDriverCommandBuffer.ConsumerHeartbeatTime()
}

func (driver *Proxy) AddSubscription(channel string, streamId int32) int64 {

	correlationId := driver.toDriverCommandBuffer.NextCorrelationId()

	logger.Debugf("driver.AddSubscription: correlationId=%d", correlationId)

	filler := func(buffer *buffer.Atomic, length *int) int32 {

		var message command.SubscriptionMessage
		message.Wrap(buffer, 0)

		message.ClientId.Set(driver.clientId)
		message.CorrelationId.Set(correlationId)
		message.RegistrationCorrelationId.Set(-1)
		message.StreamId.Set(streamId)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.ADD_SUBSCRIPTION
	}

	driver.writeCommandToDriver(filler)

	return correlationId

}

func (driver *Proxy) RemoveSubscription(registrationId int64) int64 {
	correlationId := driver.toDriverCommandBuffer.NextCorrelationId()

	logger.Debugf("driver.RemoveSubscription: correlationId=%d (subId=%d)", correlationId, registrationId)

	filler := func(buffer *buffer.Atomic, length *int) int32 {

		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.CorrelationId.Set(driver.clientId)
		message.CorrelationId.Set(correlationId)
		message.RegistrationId.Set(registrationId)

		*length = message.Size()

		return command.REMOVE_SUBSCRIPTION
	}

	driver.writeCommandToDriver(filler)

	return correlationId
}

func (driver *Proxy) AddPublication(channel string, streamId int32) int64 {

	correlationId := driver.toDriverCommandBuffer.NextCorrelationId()

	logger.Debugf("driver.AddPublication: correlationId=%d", correlationId)

	filler := func(buffer *buffer.Atomic, length *int) int32 {

		var message command.PublicationMessage
		message.Wrap(buffer, 0)
		message.ClientId.Set(driver.clientId)
		message.CorrelationId.Set(correlationId)
		message.StreamId.Set(streamId)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.ADD_PUBLICATION
	}

	driver.writeCommandToDriver(filler)

	return correlationId
}

func (driver *Proxy) RemovePublication(registrationId int64) int64 {
	correlationId := driver.toDriverCommandBuffer.NextCorrelationId()

	logger.Debugf("driver.RemovePublication: correlationId=%d (pudId=%d)", correlationId, registrationId)

	filler := func(buffer *buffer.Atomic, length *int) int32 {

		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.CorrelationId.Set(driver.clientId)
		message.CorrelationId.Set(correlationId)
		message.RegistrationId.Set(registrationId)

		*length = message.Size()

		return command.REMOVE_PUBLICATION
	}

	driver.writeCommandToDriver(filler)

	return correlationId
}

func (driver *Proxy) SendClientKeepalive() {

	filler := func(buffer *buffer.Atomic, length *int) int32 {

		var message command.CorrelatedMessage
		message.Wrap(buffer, 0)
		message.ClientId.Set(driver.clientId)
		message.CorrelationId.Set(0)

		*length = message.Size()

		return command.CLIENT_KEEPALIVE
	}

	driver.writeCommandToDriver(filler)
}

func (driver *Proxy) writeCommandToDriver(filler func(*buffer.Atomic, *int) int32) {
	messageBuffer := make([]byte, 512)

	buffer := buffer.MakeAtomic(messageBuffer)

	length := len(messageBuffer)

	msgTypeId := filler(buffer, &length)

	//fmt.Printf("DriverProxy.writeCommandToDriver: ")
	//for i := 0; i < int(length); i++ {
	//	fmt.Printf("%x ", messageBuffer[i])
	//}
	//fmt.Printf("\n")

	if !driver.toDriverCommandBuffer.Write(int32(msgTypeId), buffer, 0, int32(length)) {
		panic("couldn't write command to driver")
	}
}
