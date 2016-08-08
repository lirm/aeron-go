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
	"github.com/lirm/aeron-go/aeron/flyweight"
)

type subscriberPositionFly struct {
	flyweight.FWBase

	indicatorID    flyweight.Int32Field
	registrationID flyweight.Int64Field
}

func (m *subscriberPositionFly) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.indicatorID.Wrap(buf, pos)
	pos += m.registrationID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

type imageReadyTrailer struct {
	flyweight.FWBase

	logFile        flyweight.StringField
	sourceIdentity flyweight.StringField
}

func (m *imageReadyTrailer) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.logFile.Wrap(buf, pos, m)
	pos += m.sourceIdentity.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type errorMessage struct {
	flyweight.FWBase

	offendingCommandCorrelationID flyweight.Int64Field
	errorCode                     flyweight.Int32Field
	errorMessage                  flyweight.StringField
}

func (m *errorMessage) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.offendingCommandCorrelationID.Wrap(buf, pos)
	pos += m.errorCode.Wrap(buf, pos)
	pos += m.errorMessage.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type publicationReady struct {
	flyweight.FWBase

	correlationID          flyweight.Int64Field
	sessionID              flyweight.Int32Field
	streamID               flyweight.Int32Field
	publicationLimitOffset flyweight.Int32Field
	logFile                flyweight.StringField
}

func (m *publicationReady) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.sessionID.Wrap(buf, pos)
	pos += m.streamID.Wrap(buf, pos)
	pos += m.publicationLimitOffset.Wrap(buf, pos)
	pos += m.logFile.Wrap(buf, pos, m)

	m.SetSize(pos - offset)
	return m
}

type imageReadyHeader struct {
	flyweight.FWBase

	correlationID   flyweight.Int64Field
	sessionID       flyweight.Int32Field
	streamID        flyweight.Int32Field
	subsPosBlockLen flyweight.Int32Field
	subsPosBlockCnt flyweight.Int32Field
}

func (m *imageReadyHeader) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.sessionID.Wrap(buf, pos)
	pos += m.streamID.Wrap(buf, pos)
	pos += m.subsPosBlockLen.Wrap(buf, pos)
	pos += m.subsPosBlockCnt.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}
