package driver

import (
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/buffers"
	"log"
)

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
		sourceIdentity string, subscriberPositionCount int32,
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

func NewAdapter(driverListener    Listener, broadcastReceiver *broadcast.CopyReceiver) *ListenerAdapter {
	adapter := new(ListenerAdapter)
	adapter.listener = driverListener
	adapter.broadcastReceiver = broadcastReceiver

	return adapter
}

func (adapter *ListenerAdapter) ReceiveMessages() int {
	handler := func(msgTypeId int32, buffer *buffers.Atomic, offset int32, length int32) {
		// log.Printf("received %d", msgTypeId)
		switch int32(msgTypeId) {
		case Events.ON_PUBLICATION_READY:
			// log.Printf("received ON_PUBLICATION_READY")

			correlationId := buffer.GetInt64(offset)
			sessionId := buffer.GetInt32(offset + 8)
			streamId := buffer.GetInt32(offset + 12)
			positionLimitCounterId := buffer.GetInt32(offset + 16)
			logFileNameSize := buffer.GetInt32(offset + 20)
			logFileNameBytes := buffer.GetBytesArray(offset+24, int32(logFileNameSize))
			logFileName := string(logFileNameBytes)

			adapter.listener.OnNewPublication(streamId, sessionId, positionLimitCounterId, logFileName, correlationId)
		case Events.ON_AVAILABLE_IMAGE:
			// log.Printf("received ON_AVAILABLE_IMAGE")

			correlationId := buffer.GetInt64(offset)
			sessionId := buffer.GetInt32(offset + 8)
			streamId := buffer.GetInt32(offset + 12)
			subscriberPositionBlockLength := buffer.GetInt32(offset + 16)
			subscriberPositionCount := buffer.GetInt32(offset + 20)

			subscriberPositions := make([]SubscriberPosition, subscriberPositionCount)
			ix := 0
			pos := offset + int32(24)
			logFileNameOffset := pos + int32(subscriberPositionBlockLength)
			for ; pos < logFileNameOffset; ix++ {
				subscriberPositions[ix].indicatorId = buffer.GetInt32(pos)
				pos += 4
				subscriberPositions[ix].registrationId = buffer.GetInt64(pos)
				pos += 8
			}

			logFileNameSize := buffer.GetInt32(logFileNameOffset)
			logFileNameBytes := buffer.GetBytesArray(logFileNameOffset+4, int32(logFileNameSize))
			logFileName := string(logFileNameBytes)

			sourceIdentityOffset := logFileNameOffset + int32(4+logFileNameSize)
			sourceIdentitySize := buffer.GetInt32(sourceIdentityOffset)
			sourceIdentityBytes := buffer.GetBytesArray(sourceIdentityOffset+4, int32(sourceIdentitySize))
			sourceIdentity := string(sourceIdentityBytes)

			adapter.listener.OnAvailableImage(streamId, sessionId, logFileName,
				sourceIdentity, subscriberPositionCount, subscriberPositions, correlationId)
		case Events.ON_OPERATION_SUCCESS:
			// log.Printf("received ON_OPERATION_SUCCESS")

			correlationId := buffer.GetInt64(offset + 8)

			adapter.listener.OnOperationSuccess(correlationId)
		case Events.ON_UNAVAILABLE_IMAGE:
			// log.Printf("received ON_UNAVAILABLE_IMAGE")

			streamId := buffer.GetInt32(offset + 8)
			correlationId := buffer.GetInt64(offset)

			adapter.listener.OnUnavailableImage(streamId, correlationId)
		case Events.ON_ERROR:
			// log.Printf("received ON_ERROR")

			offendingCommandCorrelationId := buffer.GetInt64(offset)
			errorCode := buffer.GetInt32(offset + 8)
			errorSize := buffer.GetInt32(offset + 12)
			errorBytes := buffer.GetBytesArray(offset+16, int32(errorSize))
			errorMessage := string(errorBytes)

			adapter.listener.OnErrorResponse(offendingCommandCorrelationId, errorCode, errorMessage)
		default:
			log.Printf("received unhandled %d", msgTypeId)
		}
	}

	return adapter.broadcastReceiver.Receive(handler)
}
