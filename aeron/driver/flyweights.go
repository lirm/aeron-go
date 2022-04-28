/*
Copyright 2016-2018 Stanislav Liberman

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

/**
 * Control message flyweight for any errors sent from driver to clients
 *
 * <p>
 * 0                   1                   2                   3
 * 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |              Offending Command Correlation ID                 |
 * |                                                               |
 * +---------------------------------------------------------------+
 * |                         Error Code                            |
 * +---------------------------------------------------------------+
 * |                   Error Message Length                        |
 * +---------------------------------------------------------------+
 * |                       Error Message                          ...
 * ...                                                             |
 * +---------------------------------------------------------------+
 */
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
	pos += m.errorMessage.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

/**
* Message to denote that new buffers have been added for a publication.
*
* @see ControlProtocolEvents
*
* 0                   1                   2                   3
* 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
* +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
* |                         Correlation ID                        |
* |                                                               |
* +---------------------------------------------------------------+
* |                        Registration ID                        |
* |                                                               |
* +---------------------------------------------------------------+
* |                          Session ID                           |
* +---------------------------------------------------------------+
* |                           Stream ID                           |
* +---------------------------------------------------------------+
* |                   Position Limit Counter Id                   |
* +---------------------------------------------------------------+
* |                  Channel Status Indicator ID                  |
* +---------------------------------------------------------------+
* |                         Log File Length                       |
* +---------------------------------------------------------------+
* |                          Log File Name                      ...
* ...                                                             |
* +---------------------------------------------------------------+
 */
type publicationReady struct {
	flyweight.FWBase

	correlationID            flyweight.Int64Field
	registrationID           flyweight.Int64Field
	sessionID                flyweight.Int32Field
	streamID                 flyweight.Int32Field
	publicationLimitOffset   flyweight.Int32Field
	channelStatusIndicatorID flyweight.Int32Field
	logFileName              flyweight.StringField
}

func (m *publicationReady) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.registrationID.Wrap(buf, pos)
	pos += m.sessionID.Wrap(buf, pos)
	pos += m.streamID.Wrap(buf, pos)
	pos += m.publicationLimitOffset.Wrap(buf, pos)
	pos += m.channelStatusIndicatorID.Wrap(buf, pos)
	pos += m.logFileName.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

/**
 * Message to denote that a Subscription has been successfully set up.
 *
 *   0                   1                   2                   3
 *   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |                         Correlation ID                        |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                  Channel Status Indicator ID                  |
 *  +---------------------------------------------------------------+
 */
type subscriptionReady struct {
	flyweight.FWBase

	correlationID            flyweight.Int64Field
	channelStatusIndicatorID flyweight.Int32Field
}

func (m *subscriptionReady) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.channelStatusIndicatorID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

/**
* Message to denote that new buffers have been added for a subscription.
*
* NOTE: Layout should be SBE compliant
*
* @see ControlProtocolEvents
*
* 0                   1                   2                   3
* 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
* +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
* |                       Correlation ID                          |
* |                                                               |
* +---------------------------------------------------------------+
* |                         Session ID                            |
* +---------------------------------------------------------------+
* |                         Stream ID                             |
* +---------------------------------------------------------------+
* |                  Subscriber Registration Id                   |
* |                                                               |
* +---------------------------------------------------------------+
* |                    Subscriber Position Id                     |
* +---------------------------------------------------------------+
* |                       Log File Length                         |
* +---------------------------------------------------------------+
* |                        Log File Name                         ...
*...                                                              |
* +---------------------------------------------------------------+
* |                    Source identity Length                     |
* +---------------------------------------------------------------+
* |                    Source identity Name                      ...
*...                                                              |
* +---------------------------------------------------------------+
 */
type imageReadyHeader struct {
	flyweight.FWBase

	correlationID      flyweight.Int64Field
	sessionID          flyweight.Int32Field
	streamID           flyweight.Int32Field
	subsRegistrationID flyweight.Int64Field
	subsPosID          flyweight.Int32Field
	logFile            flyweight.StringField
	sourceIdentity     flyweight.StringField
}

func (m *imageReadyHeader) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.sessionID.Wrap(buf, pos)
	pos += m.streamID.Wrap(buf, pos)
	pos += m.subsRegistrationID.Wrap(buf, pos)
	pos += m.subsPosID.Wrap(buf, pos)
	pos += m.logFile.Wrap(buf, pos, m, true)
	pos += m.sourceIdentity.Wrap(buf, pos, m, true)

	m.SetSize(pos - offset)
	return m
}

/**
 * Message to denote that a Counter has been successfully set up or removed.
 *
 *   0                   1                   2                   3
 *   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 *  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 *  |                         Correlation ID                        |
 *  |                                                               |
 *  +---------------------------------------------------------------+
 *  |                           Counter ID                          |
 *  +---------------------------------------------------------------+
 */
type counterUpdate struct {
	flyweight.FWBase

	correlationID flyweight.Int64Field
	counterID     flyweight.Int32Field
}

func (m *counterUpdate) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.correlationID.Wrap(buf, pos)
	pos += m.counterID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}

/**
 * Indicate a client has timed out by the driver.
 *
 * 0                   1                   2                   3
 * 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 * |                         Client Id                             |
 * |                                                               |
 * +---------------------------------------------------------------+
 */
type clientTimeout struct {
	flyweight.FWBase

	clientID flyweight.Int64Field
}

func (m *clientTimeout) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.clientID.Wrap(buf, pos)

	m.SetSize(pos - offset)
	return m
}
