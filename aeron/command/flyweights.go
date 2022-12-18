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

package command

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/flyweight"
)

type CorrelatedMessage struct {
	flyweight.FWBase

	ClientID      flyweight.Int64Field
	CorrelationID flyweight.Int64Field
}

func (m *CorrelatedMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

type ImageMessage struct {
	flyweight.FWBase

	CorrelationID              flyweight.Int64Field
	SubscriptionRegistrationID flyweight.Int64Field
	StreamID                   flyweight.Int32Field
	Channel                    flyweight.StringField
}

func (m *ImageMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.SubscriptionRegistrationID.Wrap(buf, pos)
	pos += m.StreamID.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

type PublicationMessage struct {
	flyweight.FWBase

	ClientID      flyweight.Int64Field
	CorrelationID flyweight.Int64Field
	StreamID      flyweight.Int32Field
	Channel       flyweight.StringField
}

func (m *PublicationMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.StreamID.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

type SubscriptionMessage struct {
	flyweight.FWBase

	ClientID                  flyweight.Int64Field
	CorrelationID             flyweight.Int64Field
	RegistrationCorrelationID flyweight.Int64Field
	StreamID                  flyweight.Int32Field
	Channel                   flyweight.StringField
}

func (m *SubscriptionMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.RegistrationCorrelationID.Wrap(buf, pos)
	pos += m.StreamID.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

type RemoveMessage struct {
	flyweight.FWBase

	ClientID       flyweight.Int64Field
	CorrelationID  flyweight.Int64Field
	RegistrationID flyweight.Int64Field
}

func (m *RemoveMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.RegistrationID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

type DestinationMessage struct {
	flyweight.FWBase

	ClientID                  flyweight.Int64Field
	CorrelationID             flyweight.Int64Field
	RegistrationCorrelationID flyweight.Int64Field
	Channel                   flyweight.StringField
}

func (m *DestinationMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.RegistrationCorrelationID.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

// CounterMessage has to violate the pattern above.  The existing pattern only works either without any variable fields,
// or with exactly one variable field at the end of the message.  CounterMessage has 2 variable fields, requiring
// realignment of the second anytime there are changes to the first.  This pattern is somewhere between the Go pattern
// and the Java/C++ impls' pattern.  At some point, we may want to refactor all the flyweights to match the Java/C++
// pattern.
type CounterMessage struct {
	flyweight.FWBase

	buf           *atomic.Buffer
	offset        int
	ClientID      flyweight.Int64Field
	CorrelationID flyweight.Int64Field
	CounterTypeID flyweight.Int32Field
	key           flyweight.LengthAndRawDataField
	label         flyweight.LengthAndRawDataField
}

func (m *CounterMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	m.buf = buf
	m.offset = offset
	pos := offset
	pos += m.ClientID.Wrap(buf, pos)
	pos += m.CorrelationID.Wrap(buf, pos)
	pos += m.CounterTypeID.Wrap(buf, pos)
	pos += m.key.Wrap(buf, pos)
	m.wrapLabel()
	m.setSize()

	return m
}

func (m *CounterMessage) getLabelOffset() int {
	// Size of the first 3 fixed fields, the length field of the key, and the key buffer itself.
	return 8 + 8 + 4 + 4 + int(m.key.Length())
}

func (m *CounterMessage) setSize() {
	m.SetSize(m.getLabelOffset() + 4 + int(m.label.Length()))
}

func (m *CounterMessage) wrapLabel() {
	m.label.Wrap(m.buf, m.offset+m.getLabelOffset())
}

// Note: If you call this with a buffer that's a different length than the prior buffer, and you don't also call
// CopyLabelBuffer, the underlying data will be corrupt.  This matches the Java impl.
func (m *CounterMessage) CopyKeyBuffer(buffer *atomic.Buffer, offset int32, length int32) {
	m.key.CopyBuffer(buffer, offset, length)
}

func (m *CounterMessage) CopyLabelBuffer(buffer *atomic.Buffer, offset int32, length int32) {
	m.wrapLabel()
	m.label.CopyBuffer(buffer, offset, length)
	m.setSize()
}
