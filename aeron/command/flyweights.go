package command

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/flyweight"
)

type CorrelatedMessage struct {
	ClientId      flyweight.Int64Field
	CorrelationId flyweight.Int64Field
	length        int
}

func (m *CorrelatedMessage) Length() int {
	return m.length
}

func (m *CorrelatedMessage) Wrap(buf *buffers.Atomic, offset int) *CorrelatedMessage {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)

	m.length = pos
	return m
}

type ImageMessage struct {
	CorrelationId flyweight.Int64Field
	StreamId      flyweight.Int32Field
	Channel       flyweight.StringField
	length        int
}

func (m *ImageMessage) Length() int {
	return m.length
}

func (m *ImageMessage) Wrap(buf *buffers.Atomic, offset int) *ImageMessage {
	pos := offset
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type PublicationMessage struct {
	ClientId      flyweight.Int64Field
	CorrelationId flyweight.Int64Field
	StreamId      flyweight.Int32Field
	Channel       flyweight.StringField
	length        int
}

func (m *PublicationMessage) Length() int {
	return m.length
}

func (m *PublicationMessage) Wrap(buf *buffers.Atomic, offset int) *PublicationMessage {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type SubscriptionMessage struct {
	ClientId                  flyweight.Int64Field
	CorrelationId             flyweight.Int64Field
	RegistrationCorrelationId flyweight.Int64Field
	StreamId                  flyweight.Int32Field
	Channel                   flyweight.StringField
	length                    int
}

func (m *SubscriptionMessage) Length() int {
	return m.length
}

func (m *SubscriptionMessage) Wrap(buf *buffers.Atomic, offset int) *SubscriptionMessage {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.RegistrationCorrelationId.Wrap(buf, pos)
	pos += m.StreamId.Wrap(buf, pos)
	pos += m.Channel.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type RemoveMessage struct {
	ClientId       flyweight.Int64Field
	CorrelationId  flyweight.Int64Field
	RegistrationId flyweight.Int64Field
	length         int
}

func (m *RemoveMessage) Length() int {
	return m.length
}

func (m *RemoveMessage) Wrap(buf *buffers.Atomic, offset int) *RemoveMessage {
	pos := offset
	pos += m.ClientId.Wrap(buf, pos)
	pos += m.CorrelationId.Wrap(buf, pos)
	pos += m.RegistrationId.Wrap(buf, pos)

	m.length = pos - offset
	return m
}
