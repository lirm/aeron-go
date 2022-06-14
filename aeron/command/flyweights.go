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
