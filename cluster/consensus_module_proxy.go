package cluster

import (
	"fmt"

	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/cluster/codecs"
)

const (
	scheduleTimerBlockLength = 16
	cancelTimerBlockLength   = 8
)

// Proxy class for encapsulating encoding and sending of control protocol messages to a cluster
type consensusModuleProxy struct {
	marshaller    *codecs.SbeGoMarshaller // currently shared as we're not reentrant (but could be here)
	idleStrategy  idlestrategy.Idler
	rangeChecking bool
	publication   *aeron.Publication
	buffer        *atomic.Buffer
}

func newConsensusModuleProxy(
	options *Options,
	publication *aeron.Publication,
) *consensusModuleProxy {
	return &consensusModuleProxy{
		marshaller:    codecs.NewSbeGoMarshaller(),
		rangeChecking: options.RangeChecking,
		publication:   publication,
		buffer:        atomic.MakeBuffer(make([]byte, 500)),
	}
}

// From here we have all the functions that create a data packet and send it on the
// publication. Responses will be processed on the control

// ConnectRequest packet and send
func (proxy *consensusModuleProxy) serviceAckRequest(
	logPosition int64,
	timestamp int64,
	ackID int64,
	relevantID int64,
	serviceID int32,
) {
	// Create a packet and send it
	bytes, err := codecs.ServiceAckRequestPacket(
		proxy.marshaller,
		proxy.rangeChecking,
		logPosition,
		timestamp,
		ackID,
		relevantID,
		serviceID,
	)
	if err != nil {
		panic(err)
	}
	proxy.send(bytes)
}

func (proxy *consensusModuleProxy) closeSessionRequest(
	clusterSessionId int64,
) {
	// Create a packet and send it
	bytes, err := codecs.CloseSessionRequestPacket(
		proxy.marshaller,
		proxy.rangeChecking,
		clusterSessionId,
	)
	if err != nil {
		panic(err)
	}
	proxy.send(bytes)
}

func (proxy *consensusModuleProxy) scheduleTimer(correlationId int64, deadline int64) bool {
	buf := proxy.initBuffer(scheduleTimerTemplateId, scheduleTimerBlockLength)
	buf.PutInt64(SBEHeaderLength, correlationId)
	buf.PutInt64(SBEHeaderLength+8, deadline)
	return proxy.offer(buf, SBEHeaderLength+scheduleTimerBlockLength) >= 0
}

func (proxy *consensusModuleProxy) cancelTimer(correlationId int64) bool {
	buf := proxy.initBuffer(cancelTimerTemplateId, cancelTimerBlockLength)
	buf.PutInt64(SBEHeaderLength, correlationId)
	return proxy.offer(buf, SBEHeaderLength+cancelTimerBlockLength) >= 0
}

func (proxy *consensusModuleProxy) initBuffer(templateId uint16, blockLength uint16) *atomic.Buffer {
	buf := proxy.buffer
	buf.PutUInt16(0, blockLength)
	buf.PutUInt16(2, templateId)
	buf.PutUInt16(4, clusterSchemaId)
	buf.PutUInt16(6, clusterSchemaVersion)
	return buf
}

// send to our request publication
func (proxy *consensusModuleProxy) send(payload []byte) {
	buffer := atomic.MakeBuffer(payload)
	for proxy.offer(buffer, buffer.Capacity()) < 0 {
		proxy.idleStrategy.Idle(0)
	}
}

func (proxy *consensusModuleProxy) offer(buffer *atomic.Buffer, length int32) int64 {
	var result int64
	for i := 0; i < 3; i++ {
		result = proxy.publication.Offer(buffer, 0, length, nil)
		if result >= 0 {
			break
		} else if result == aeron.NotConnected || result == aeron.PublicationClosed || result == aeron.MaxPositionExceeded {
			panic(fmt.Sprintf("offer failed, result=%d", result))
		}
	}
	return result
}
