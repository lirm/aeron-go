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
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/flyweight"
)

type SubscriberPositionFly struct {
	indicatorId    flyweight.Int32Field
	registrationId flyweight.Int64Field
	length         int
}

func (m *SubscriberPositionFly) Length() int {
	return m.length
}

func (m *SubscriberPositionFly) Wrap(buf *buffers.Atomic, offset int) *SubscriberPositionFly {
	pos := offset
	pos += m.indicatorId.Wrap(buf, pos)
	pos += m.registrationId.Wrap(buf, pos)

	m.length = pos - offset
	return m
}

type ImageReadyTrailer struct {
	logFile        flyweight.StringField
	sourceIdentity flyweight.StringField
	length         int
}

func (m *ImageReadyTrailer) Length() int {
	return m.length
}

func (m *ImageReadyTrailer) Wrap(buf *buffers.Atomic, offset int) *ImageReadyTrailer {
	pos := offset
	pos += m.logFile.Wrap(buf, pos, &m.length)
	pos += m.sourceIdentity.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type ErrorMessage struct {
	offendingCommandCorrelationId flyweight.Int64Field
	errorCode                     flyweight.Int32Field
	errorMessage                  flyweight.StringField
	length                        int
}

func (m *ErrorMessage) Length() int {
	return m.length
}

func (m *ErrorMessage) Wrap(buf *buffers.Atomic, offset int) *ErrorMessage {
	pos := offset
	pos += m.offendingCommandCorrelationId.Wrap(buf, pos)
	pos += m.errorCode.Wrap(buf, pos)
	pos += m.errorMessage.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type PublicationReady struct {
	correlationId          flyweight.Int64Field
	sessionId              flyweight.Int32Field
	streamId               flyweight.Int32Field
	publicationLimitOffset flyweight.Int32Field
	logFile                flyweight.StringField
	length                 int
}

func (m *PublicationReady) Length() int {
	return m.length
}

func (m *PublicationReady) Wrap(buf *buffers.Atomic, offset int) *PublicationReady {
	pos := offset
	pos += m.correlationId.Wrap(buf, pos)
	pos += m.sessionId.Wrap(buf, pos)
	pos += m.streamId.Wrap(buf, pos)
	pos += m.publicationLimitOffset.Wrap(buf, pos)
	pos += m.logFile.Wrap(buf, pos, &m.length)

	m.length = pos - offset
	return m
}

type ImageReadyHeader struct {
	correlationId   flyweight.Int64Field
	sessionId       flyweight.Int32Field
	streamId        flyweight.Int32Field
	subsPosBlockLen flyweight.Int32Field
	subsPosBlockCnt flyweight.Int32Field
	length          int
}

func (m *ImageReadyHeader) Length() int {
	return m.length
}

func (m *ImageReadyHeader) Wrap(buf *buffers.Atomic, offset int) *ImageReadyHeader {
	pos := offset
	pos += m.correlationId.Wrap(buf, pos)
	pos += m.sessionId.Wrap(buf, pos)
	pos += m.streamId.Wrap(buf, pos)
	pos += m.subsPosBlockLen.Wrap(buf, pos)
	pos += m.subsPosBlockCnt.Wrap(buf, pos)

	m.length = pos - offset
	return m
}
