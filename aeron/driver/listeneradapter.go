package driver

import (
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/command"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("driver")

var Events = struct {
	/** Error Response */
	ON_ERROR int32
	/** New subscription Buffer Notification */
	ON_AVAILABLE_IMAGE int32
	/** New publication Buffer Notification */
	ON_PUBLICATION_READY int32
	/** Operation Succeeded */
	ON_OPERATION_SUCCESS int32

	/** Inform client of timeout and removal of inactive image */
	ON_UNAVAILABLE_IMAGE int32
}{
	0x0F01,
	0x0F02,
	0x0F03,
	0x0F04,

	0x0F05,
}

type SubscriberPosition struct {
	indicatorId    int32
	registrationId int64
}

func (pos *SubscriberPosition) RegistrationId() int64 {
	return pos.registrationId
}

func (pos *SubscriberPosition) IndicatorId() int32 {
	return pos.indicatorId
}

type Listener interface {
	OnNewPublication(streamId int32, sessionId int32, positionLimitCounterId int32,
		logFileName string, registrationId int64)
	OnAvailableImage(streamId int32, sessionId int32, logFilename string,
		sourceIdentity string, subscriberPositionCount int,
		subscriberPositions []SubscriberPosition,
		correlationId int64)
	OnUnavailableImage(streamId int32, correlationId int64)
	OnOperationSuccess(correlationId int64)
	OnErrorResponse(offendingCommandCorrelationId int64, errorCode int32, errorMessage string)
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
	handler := func(msgTypeId int32, buffer *buffers.Atomic, offset int32, length int32) {
		logger.Debugf("received %d", msgTypeId)
		switch int32(msgTypeId) {
		case Events.ON_PUBLICATION_READY:
			logger.Debugf("received ON_PUBLICATION_READY")

			var msg PublicationReady
			msg.Wrap(buffer, int(offset))

			correlationId := msg.correlationId.Get()
			sessionId := msg.sessionId.Get()
			streamId := msg.streamId.Get()
			positionLimitCounterId := msg.publicationLimitOffset.Get()
			logFileName := msg.logFile.Get()

			adapter.listener.OnNewPublication(streamId, sessionId, positionLimitCounterId, logFileName, correlationId)
		case Events.ON_AVAILABLE_IMAGE:
			logger.Debugf("received ON_AVAILABLE_IMAGE")

			var header ImageReadyHeader
			header.Wrap(buffer, int(offset))

			correlationId := header.correlationId.Get()
			sessionId := header.sessionId.Get()
			streamId := header.streamId.Get()
			subsPosBlockLen := header.subsPosBlockLen.Get()
			subsPosBlockCnt := int(header.subsPosBlockCnt.Get())
			logger.Debugf("position count: %d block len: %d", subsPosBlockCnt, subsPosBlockLen)

			subscriberPositions := make([]SubscriberPosition, subsPosBlockCnt)
			pos := offset + int32(24)
			var posFly SubscriberPositionFly
			for ix := 0; ix < subsPosBlockCnt; ix++ {
				posFly.Wrap(buffer, int(pos))
				pos += subsPosBlockLen

				subscriberPositions[ix].indicatorId = posFly.indicatorId.Get()
				subscriberPositions[ix].registrationId = posFly.registrationId.Get()
			}
			logger.Debugf("positions: %v", subscriberPositions)

			var trailer ImageReadyTrailer
			trailer.Wrap(buffer, int(pos))

			logFileName := trailer.logFile.Get()
			logger.Debugf("logFileName: %v", logFileName)

			sourceIdentity := trailer.sourceIdentity.Get()
			logger.Debugf("sourceIdentity: %v", sourceIdentity)

			adapter.listener.OnAvailableImage(streamId, sessionId, logFileName, sourceIdentity,
				subsPosBlockCnt, subscriberPositions, correlationId)
		case Events.ON_OPERATION_SUCCESS:
			logger.Debugf("received ON_OPERATION_SUCCESS")

			var msg command.CorrelatedMessage
			msg.Wrap(buffer, int(offset))

			correlationId := msg.CorrelationId.Get()

			adapter.listener.OnOperationSuccess(correlationId)
		case Events.ON_UNAVAILABLE_IMAGE:
			logger.Debugf("received ON_UNAVAILABLE_IMAGE")

			var msg command.ImageMessage
			msg.Wrap(buffer, int(offset))

			streamId := msg.StreamId.Get()
			correlationId := msg.CorrelationId.Get()

			adapter.listener.OnUnavailableImage(streamId, correlationId)
		case Events.ON_ERROR:
			logger.Debugf("received ON_ERROR")

			var msg ErrorMessage
			msg.Wrap(buffer, int(offset))

			adapter.listener.OnErrorResponse(msg.offendingCommandCorrelationId.Get(),
				msg.errorCode.Get(), msg.errorMessage.Get())
		default:
			logger.Debugf("received unhandled %d", msgTypeId)
		}
	}

	return adapter.broadcastReceiver.Receive(handler)
}
