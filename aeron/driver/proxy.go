/*
Copyright 2016 Stanislav Liberman
Copyright (C) 2022 Talos, Inc.

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
	"errors"
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
	logger.Infof("aeron clientID:%d", driver.clientID)
	return driver
}

// ClientID returns the client ID for this connection to the driver.
func (driver *Proxy) ClientID() int64 {
	return driver.clientID
}

// TimeOfLastDriverKeepalive gets the time of the last keep alive update sent to media driver
func (driver *Proxy) TimeOfLastDriverKeepalive() int64 {
	return driver.toDriverCommandBuffer.ConsumerHeartbeatTime()
}

// NextCorrelationID generates the next correlation id that is unique for the connected Media Driver.
func (driver *Proxy) NextCorrelationID() int64 {
	return driver.toDriverCommandBuffer.NextCorrelationID()
}

// AddSubscription sends driver command to add new subscription
func (driver *Proxy) AddSubscription(channel string, streamID int32) (int64, error) {

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

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// RemoveSubscription sends driver command to remove subscription
func (driver *Proxy) RemoveSubscription(registrationID int64) error {
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

	return driver.writeCommandToDriver(filler)
}

// AddPublication sends driver command to add new publication
func (driver *Proxy) AddPublication(channel string, streamID int32) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddPublication: clientId=%d correlationId=%d",
		driver.clientID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.PublicationMessage
		message.Wrap(buffer, 0)
		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.StreamID.Set(streamID)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.AddPublication
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// AddExclusivePublication sends driver command to add new publication
func (driver *Proxy) AddExclusivePublication(channel string, streamID int32) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddExclusivePublication: clientId=%d correlationId=%d",
		driver.clientID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.PublicationMessage
		message.Wrap(buffer, 0)
		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.StreamID.Set(streamID)
		message.Channel.Set(channel)

		*length = message.Size()

		return command.AddExclusivePublication
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// RemovePublication sends driver command to remove publication
func (driver *Proxy) RemovePublication(registrationID int64) error {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemovePublication: clientId=%d correlationId=%d (regId=%d)",
		driver.clientID, correlationID, registrationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.RegistrationID.Set(registrationID)

		*length = message.Size()

		return command.RemovePublication
	}

	return driver.writeCommandToDriver(filler)
}

// ClientClose sends a client close to the driver.
func (driver *Proxy) ClientClose() error {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.ClientClose: clientId=%d correlationId=%d",
		driver.clientID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.CorrelatedMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)

		*length = message.Size()

		return command.ClientClose
	}

	return driver.writeCommandToDriver(filler)
}

// AddDestination sends driver command to add a destination to an existing Publication.
func (driver *Proxy) AddDestination(registrationID int64, channel string) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddDestination: clientID=%d registrationID=%d correlationID=%d",
		driver.clientID, registrationID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.DestinationMessage
		message.Wrap(buffer, 0)
		message.RegistrationCorrelationID.Set(registrationID)
		message.Channel.Set(channel)
		message.CorrelationID.Set(correlationID)
		message.ClientID.Set(driver.clientID)

		*length = message.Size()

		return command.AddDestination
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// RemoveDestination sends driver command to remove a destination from an existing Publication.
func (driver *Proxy) RemoveDestination(registrationID int64, channel string) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemoveDestination: clientID=%d registrationID=%d correlationID=%d",
		driver.clientID, registrationID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.DestinationMessage
		message.Wrap(buffer, 0)
		message.RegistrationCorrelationID.Set(registrationID)
		message.Channel.Set(channel)
		message.CorrelationID.Set(correlationID)
		message.ClientID.Set(driver.clientID)

		*length = message.Size()

		return command.RemoveDestination
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// AddRcvDestination sends driver command to add a destination to the receive
// channel of an existing MDS Subscription.
func (driver *Proxy) AddRcvDestination(registrationID int64, channel string) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.AddRcvDestination: clientID=%d registrationID=%d correlationID=%d channel=%s",
		driver.clientID, registrationID, correlationID, channel)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.DestinationMessage
		message.Wrap(buffer, 0)
		message.RegistrationCorrelationID.Set(registrationID)
		message.Channel.Set(channel)
		message.CorrelationID.Set(correlationID)
		message.ClientID.Set(driver.clientID)

		*length = message.Size()

		return command.AddRcvDestination
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, nil
	}
}

// RemoveRcvDestination sends driver command to remove a destination from the
// receive channel of an existing MDS Subscription.
func (driver *Proxy) RemoveRcvDestination(registrationID int64, channel string) (int64, error) {

	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemoveRcvDestination: clientID=%d registrationID=%d correlationID=%d",
		driver.clientID, registrationID, correlationID)

	filler := func(buffer *atomic.Buffer, length *int) int32 {

		var message command.DestinationMessage
		message.Wrap(buffer, 0)
		message.RegistrationCorrelationID.Set(registrationID)
		message.Channel.Set(channel)
		message.CorrelationID.Set(correlationID)
		message.ClientID.Set(driver.clientID)

		*length = message.Size()

		return command.RemoveRcvDestination
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// AddCounter adds a new counter with a type id plus the label and key are provided in buffers.
func (driver *Proxy) AddCounter(typeId int32, keyBuffer *atomic.Buffer, keyOffset int32, keyLength int32,
	labelBuffer *atomic.Buffer, labelOffset int32, labelLength int32) (int64, error) {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	filler := func(buffer *atomic.Buffer, length *int) int32 {
		var message command.CounterMessage
		message.Wrap(buffer, 0)
		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.CounterTypeID.Set(typeId)
		message.CopyKeyBuffer(keyBuffer, keyOffset, keyLength)
		message.CopyLabelBuffer(labelBuffer, labelOffset, labelLength)

		*length = message.Size()

		return command.AddCounter
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// AddCounterByLabel adds a new counter with a type id and label.  The key will be blank.
func (driver *Proxy) AddCounterByLabel(typeId int32, label string) (int64, error) {
	correlationID := driver.toDriverCommandBuffer.NextCorrelationID()

	filler := func(buffer *atomic.Buffer, length *int) int32 {
		var message command.CounterMessage
		message.Wrap(buffer, 0)
		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationID)
		message.CounterTypeID.Set(typeId)
		message.CopyLabelString(label)

		*length = message.Size()

		return command.AddCounter
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationID, nil
	} else {
		return 0, err
	}
}

// RemoveCounter instructs the media driver to remove an existing counter by its registration id.  Returns the
// correlation id for the command.

// RemoveCounter instructs the media driver to remove an existing counter by its registration id.  Returns the
// correlation id for the command.
func (driver *Proxy) RemoveCounter(registrationId int64) (int64, error) {
	correlationId := driver.toDriverCommandBuffer.NextCorrelationID()

	logger.Debugf("driver.RemoveCounter: correlationId=%d (counter registrationId=%d)", correlationId, registrationId)

	filler := func(buffer *atomic.Buffer, length *int) int32 {
		var message command.RemoveMessage
		message.Wrap(buffer, 0)

		message.ClientID.Set(driver.clientID)
		message.CorrelationID.Set(correlationId)
		message.RegistrationID.Set(registrationId)

		*length = message.Size()

		return command.RemoveCounter
	}

	if err := driver.writeCommandToDriver(filler); err == nil {
		return correlationId, nil
	} else {
		return 0, err
	}
}

func (driver *Proxy) writeCommandToDriver(filler func(*atomic.Buffer, *int) int32) error {
	messageBuffer := make([]byte, 512)

	buffer := atomic.MakeBuffer(messageBuffer)

	length := len(messageBuffer)

	msgTypeID := filler(buffer, &length)

	if driver.toDriverCommandBuffer.Write(int32(msgTypeID), buffer, 0, int32(length)) {
		return nil
	} else {
		return errors.New("couldn't write command to driver")
	}
}
