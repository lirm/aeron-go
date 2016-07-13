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

func (m *CorrelatedMessage) Wrap(buf *buffers.Atomic) *CorrelatedMessage {
	offset := 0
	offset += m.ClientId.Wrap(buf, offset)
	offset += m.CorrelationId.Wrap(buf, offset)

	m.length = offset
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

func (m *ImageMessage) Wrap(buf *buffers.Atomic) *ImageMessage {
	offset := 0
	offset += m.CorrelationId.Wrap(buf, offset)
	offset += m.StreamId.Wrap(buf, offset)
	offset += m.Channel.Wrap(buf, offset, &m.length)

	m.length = offset
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

func (m *PublicationMessage) Wrap(buf *buffers.Atomic) *PublicationMessage {
	offset := 0
	offset += m.ClientId.Wrap(buf, offset)
	offset += m.CorrelationId.Wrap(buf, offset)
	offset += m.StreamId.Wrap(buf, offset)
	offset += m.Channel.Wrap(buf, offset, &m.length)

	m.length = offset
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

func (m *SubscriptionMessage) Wrap(buf *buffers.Atomic) *SubscriptionMessage {
	offset := 0
	offset += m.ClientId.Wrap(buf, offset)
	offset += m.CorrelationId.Wrap(buf, offset)
	offset += m.RegistrationCorrelationId.Wrap(buf, offset)
	offset += m.StreamId.Wrap(buf, offset)
	offset += m.Channel.Wrap(buf, offset, &m.length)

	m.length = offset
	return m
}

type RemoveMessage struct {
	ClientId      flyweight.Int64Field
	CorrelationId flyweight.Int64Field
	RegistrationId flyweight.Int64Field
	length        int
}

func (m *RemoveMessage) Length() int {
	return m.length
}

func (m *RemoveMessage) Wrap(buf *buffers.Atomic) *RemoveMessage {
	offset := 0
	offset += m.ClientId.Wrap(buf, offset)
	offset += m.CorrelationId.Wrap(buf, offset)
	offset += m.RegistrationId.Wrap(buf, offset)

	m.length = offset
	return m
}