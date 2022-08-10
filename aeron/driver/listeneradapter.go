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
	"strings"

	"github.com/lirm/aeron-go/aeron/logging"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/command"
)

var logger = logging.MustGetLogger("driver")

var Events = struct {
	/** Error Response */
	OnError int32
	/** New subscription Buffer Notification */
	OnAvailableImage int32
	/** New publication Buffer Notification */
	OnPublicationReady int32
	/** Operation Succeeded */
	OnOperationSuccess int32
	/** Inform client of timeout and removal of inactive image */
	OnUnavailableImage int32
	/** New Exclusive Publication Buffer notification */
	OnExclusivePublicationReady int32
	/** New subscription notification */
	OnSubscriptionReady int32
	/** New counter notification */
	OnCounterReady int32
	/** inform clients of removal of counter */
	OnUnavailableCounter int32
	/** inform clients of client timeout */
	OnClientTimeout int32
}{
	0x0F01,
	0x0F02,
	0x0F03,
	0x0F04,
	0x0F05,
	0x0F06,
	0x0F07,
	0x0F08,
	0x0F09,
	0x0F0A,
}

type SubscriberPosition struct {
	indicatorID    int32
	registrationID int64
}

func (pos *SubscriberPosition) RegistrationID() int64 {
	return pos.registrationID
}

func (pos *SubscriberPosition) IndicatorID() int32 {
	return pos.indicatorID
}

type Listener interface {
	OnNewPublication(streamID int32, sessionID int32, positionLimitCounterID int32, channelStatusIndicatorID int32,
		logFileName string, correlationID int64, registrationID int64)
	OnNewExclusivePublication(streamID int32, sessionID int32, positionLimitCounterID int32, channelStatusIndicatorID int32,
		logFileName string, correlationID int64, registrationID int64)
	OnAvailableImage(streamID int32, sessionID int32, logFilename string, sourceIdentity string,
		subscriberPositionID int32, subsRegID int64, correlationID int64)
	OnUnavailableImage(correlationID int64, subscriptionRegistrationID int64)
	OnOperationSuccess(correlationID int64)
	OnErrorResponse(offendingCommandCorrelationID int64, errorCode int32, errorMessage string)
	OnChannelEndpointError(correlationID int64, errorMessage string)
	OnSubscriptionReady(correlationID int64, channelStatusIndicatorID int32)
	OnAvailableCounter(correlationID int64, counterID int32)
	OnUnavailableCounter(correlationID int64, counterID int32)
	OnClientTimeout(clientID int64)
}

type ListenerAdapter struct {
	listener          Listener
	broadcastReceiver *broadcast.CopyReceiver
}

func NewAdapter(driverListener Listener, broadcastReceiver *broadcast.CopyReceiver) *ListenerAdapter {
	adapter := new(ListenerAdapter)
	adapter.listener = driverListener
	adapter.broadcastReceiver = broadcastReceiver

	return adapter
}

func (adapter *ListenerAdapter) ReceiveMessages() int {
	handler := func(msgTypeID int32, buffer *atomic.Buffer, offset int32, length int32) {
		logger.Debugf("received %d", msgTypeID)
		switch int32(msgTypeID) {
		case Events.OnPublicationReady:
			logger.Debugf("received ON_PUBLICATION_READY")

			var msg publicationReady
			msg.Wrap(buffer, int(offset))

			streamID := msg.streamID.Get()
			sessionID := msg.sessionID.Get()
			positionLimitCounterID := msg.publicationLimitOffset.Get()
			channelStatusIndicatorID := msg.channelStatusIndicatorID.Get()
			correlationID := msg.correlationID.Get()
			registrationID := msg.registrationID.Get()
			logFileName := msg.logFileName.Get()

			adapter.listener.OnNewPublication(streamID, sessionID, positionLimitCounterID, channelStatusIndicatorID,
				logFileName, correlationID, registrationID)
		case Events.OnExclusivePublicationReady:
			logger.Debugf("received ON_EXCLUSIVE_PUBLICATION_READY")

			var msg publicationReady
			msg.Wrap(buffer, int(offset))

			streamID := msg.streamID.Get()
			sessionID := msg.sessionID.Get()
			positionLimitCounterID := msg.publicationLimitOffset.Get()
			channelStatusIndicatorID := msg.channelStatusIndicatorID.Get()
			correlationID := msg.correlationID.Get()
			registrationID := msg.registrationID.Get()
			logFileName := msg.logFileName.Get()

			adapter.listener.OnNewExclusivePublication(streamID, sessionID, positionLimitCounterID, channelStatusIndicatorID,
				logFileName, correlationID, registrationID)
		case Events.OnSubscriptionReady:
			logger.Debugf("received ON_SUBSCRIPTION_READY")

			var msg subscriptionReady
			msg.Wrap(buffer, int(offset))

			correlationID := msg.correlationID.Get()
			channelStatusIndicatorID := msg.channelStatusIndicatorID.Get()

			adapter.listener.OnSubscriptionReady(correlationID, channelStatusIndicatorID)
		case Events.OnAvailableImage:
			logger.Debugf("received ON_AVAILABLE_IMAGE")

			var header imageReadyHeader
			header.Wrap(buffer, int(offset))

			streamID := header.streamID.Get()
			sessionID := header.sessionID.Get()
			logFileName := header.logFile.Get()
			sourceIdentity := header.sourceIdentity.Get()
			subsPosID := header.subsPosID.Get()
			subsRegID := header.subsRegistrationID.Get()
			correlationID := header.correlationID.Get()

			logger.Debugf("logFileName: %v", logFileName)
			logger.Debugf("sourceIdentity: %v", sourceIdentity)

			adapter.listener.OnAvailableImage(streamID, sessionID, logFileName, sourceIdentity, subsPosID, subsRegID,
				correlationID)
		case Events.OnOperationSuccess:
			logger.Debugf("received ON_OPERATION_SUCCESS")

			var msg command.CorrelatedMessage
			msg.Wrap(buffer, int(offset))

			correlationID := msg.CorrelationID.Get()

			adapter.listener.OnOperationSuccess(correlationID)
		case Events.OnUnavailableImage:
			logger.Debugf("received ON_UNAVAILABLE_IMAGE")

			var msg command.ImageMessage
			msg.Wrap(buffer, int(offset))

			correlationID := msg.CorrelationID.Get()
			subscriptionRegistrationID := msg.SubscriptionRegistrationID.Get()

			adapter.listener.OnUnavailableImage(correlationID, subscriptionRegistrationID)
		case Events.OnError:
			logger.Debugf("received ON_ERROR")

			var msg errorMessage
			msg.Wrap(buffer, int(offset))

			if msg.errorCode.Get() == command.ErrorCodeChannelEndpointError ||
				strings.Contains(msg.errorMessage.Get(), "Address already in use") { // hack for c media driver
				adapter.listener.OnChannelEndpointError(msg.offendingCommandCorrelationID.Get(), msg.errorMessage.Get())
			} else {
				adapter.listener.OnErrorResponse(msg.offendingCommandCorrelationID.Get(),
					msg.errorCode.Get(), msg.errorMessage.Get())
			}
		case Events.OnCounterReady:
			logger.Debugf("received ON_COUNTER_READY")

			var msg counterUpdate
			msg.Wrap(buffer, int(offset))

			adapter.listener.OnAvailableCounter(msg.correlationID.Get(), msg.counterID.Get())
		case Events.OnUnavailableCounter:
			logger.Debugf("received ON_UNAVAILABLE_COUNTER")

			var msg counterUpdate
			msg.Wrap(buffer, int(offset))

			adapter.listener.OnUnavailableCounter(msg.correlationID.Get(), msg.counterID.Get())
		case Events.OnClientTimeout:
			logger.Debugf("received ON_CLIENT_TIMEOUT")

			var msg clientTimeout
			msg.Wrap(buffer, int(offset))

			adapter.listener.OnClientTimeout(msg.clientID.Get())
		default:
			logger.Fatalf("received unhandled %d", msgTypeID)
		}
	}

	return adapter.broadcastReceiver.Receive(handler)
}
