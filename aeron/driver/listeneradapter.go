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
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/command"
	"github.com/op/go-logging"
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
}{
	0x0F01,
	0x0F02,
	0x0F03,
	0x0F04,

	0x0F05,
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
	OnNewPublication(streamID int32, sessionID int32, positionLimitCounterID int32,
		logFileName string, registrationID int64)
	OnAvailableImage(streamID int32, sessionID int32, logFilename string,
		sourceIdentity string, subscriberPositionCount int,
		subscriberPositions []SubscriberPosition,
		correlationID int64)
	OnUnavailableImage(streamID int32, correlationID int64)
	OnOperationSuccess(correlationID int64)
	OnErrorResponse(offendingCommandCorrelationID int64, errorCode int32, errorMessage string)
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

			var msg PublicationReady
			msg.Wrap(buffer, int(offset))

			correlationID := msg.correlationID.Get()
			sessionID := msg.sessionID.Get()
			streamID := msg.streamID.Get()
			positionLimitCounterID := msg.publicationLimitOffset.Get()
			logFileName := msg.logFile.Get()

			adapter.listener.OnNewPublication(streamID, sessionID, positionLimitCounterID, logFileName, correlationID)
		case Events.OnAvailableImage:
			logger.Debugf("received ON_AVAILABLE_IMAGE")

			var header ImageReadyHeader
			header.Wrap(buffer, int(offset))

			correlationID := header.correlationID.Get()
			sessionID := header.sessionID.Get()
			streamID := header.streamID.Get()
			subsPosBlockLen := header.subsPosBlockLen.Get()
			subsPosBlockCnt := int(header.subsPosBlockCnt.Get())
			logger.Debugf("position count: %d block len: %d", subsPosBlockCnt, subsPosBlockLen)

			subscriberPositions := make([]SubscriberPosition, subsPosBlockCnt)
			pos := offset + int32(24)
			var posFly SubscriberPositionFly
			for ix := 0; ix < subsPosBlockCnt; ix++ {
				posFly.Wrap(buffer, int(pos))
				pos += subsPosBlockLen

				subscriberPositions[ix].indicatorID = posFly.indicatorID.Get()
				subscriberPositions[ix].registrationID = posFly.registrationID.Get()
			}
			logger.Debugf("positions: %v", subscriberPositions)

			var trailer ImageReadyTrailer
			trailer.Wrap(buffer, int(pos))

			logFileName := trailer.logFile.Get()
			logger.Debugf("logFileName: %v", logFileName)

			sourceIdentity := trailer.sourceIdentity.Get()
			logger.Debugf("sourceIdentity: %v", sourceIdentity)

			adapter.listener.OnAvailableImage(streamID, sessionID, logFileName, sourceIdentity,
				subsPosBlockCnt, subscriberPositions, correlationID)
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

			streamID := msg.StreamID.Get()
			correlationID := msg.CorrelationID.Get()

			adapter.listener.OnUnavailableImage(streamID, correlationID)
		case Events.OnError:
			logger.Debugf("received ON_ERROR")

			var msg ErrorMessage
			msg.Wrap(buffer, int(offset))

			adapter.listener.OnErrorResponse(msg.offendingCommandCorrelationID.Get(),
				msg.errorCode.Get(), msg.errorMessage.Get())
		default:
			logger.Debugf("received unhandled %d", msgTypeID)
		}
	}

	return adapter.broadcastReceiver.Receive(handler)
}
