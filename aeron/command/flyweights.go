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

package command

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/flyweight"
)

type CorrelatedMessage struct {
	flyweight.FWBase

	ClientId      flyweight.Int64Field
	CorrelationId flyweight.Int64Field
}

func (m *CorrelatedMessage) Wrap(buf *buffers.Atomic, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

type ImageMessage struct {
	flyweight.FWBase

	CorrelationId flyweight.Int64Field
	StreamId      flyweight.Int32Field
	Channel       flyweight.StringField
}

func (m *ImageMessage) Wrap(buf *buffers.Atomic, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type PublicationMessage struct {
	flyweight.FWBase

	ClientId      flyweight.Int64Field
	CorrelationId flyweight.Int64Field
	StreamId      flyweight.Int32Field
	Channel       flyweight.StringField
}

func (m *PublicationMessage) Wrap(buf *buffers.Atomic, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type SubscriptionMessage struct {
	flyweight.FWBase

	ClientId                  flyweight.Int64Field
	CorrelationId             flyweight.Int64Field
	RegistrationCorrelationId flyweight.Int64Field
	StreamId                  flyweight.Int32Field
	Channel                   flyweight.StringField
}

func (m *SubscriptionMessage) Wrap(buf *buffers.Atomic, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.RegistrationCorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type RemoveMessage struct {
	flyweight.FWBase

	ClientId       flyweight.Int64Field
	CorrelationId  flyweight.Int64Field
	RegistrationId flyweight.Int64Field
}

func (m *RemoveMessage) Wrap(buf *buffers.Atomic, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.RegistrationId.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}
