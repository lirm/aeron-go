package cluster

import (
	"bytes"

	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/cluster/codecs"
)

type ServiceAdapter struct {
	marshaller   *codecs.SbeGoMarshaller
	agent        *ClusteredServiceAgent
	subscription *aeron.Subscription
}

func (adapter *ServiceAdapter) Poll() int {
	if adapter.subscription.IsClosed() {
		panic("subscription closed")
	}
	return adapter.subscription.Poll(adapter.onFragment, 10)
}

func (adapter *ServiceAdapter) onFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	if length < SBEHeaderLength {
		return
	}
	blockLength := buffer.GetUInt16(offset)
	templateId := buffer.GetUInt16(offset + 2)
	schemaId := buffer.GetUInt16(offset + 4)
	version := buffer.GetUInt16(offset + 6)
	if schemaId != clusterSchemaId {
		logger.Errorf("ServiceAdapter: unexpected schemaId=%d templateId=%d blockLen=%d version=%d",
			schemaId, templateId, blockLength, version)
		return
	}
	offset += SBEHeaderLength
	length -= SBEHeaderLength

	switch templateId {
	case joinLogTemplateId:
		buf := &bytes.Buffer{}
		buffer.WriteBytes(buf, offset, length)
		joinLog := &codecs.JoinLog{}
		if err := joinLog.Decode(adapter.marshaller, buf, version, blockLength, true); err != nil {
			logger.Errorf("ServiceAdapter: join log decode error: %v", err)
		} else {
			adapter.agent.onJoinLog(
				joinLog.LogPosition,
				joinLog.MaxLogPosition,
				joinLog.MemberId,
				joinLog.LogSessionId,
				joinLog.LogStreamId,
				joinLog.IsStartup == codecs.BooleanType.TRUE,
				Role(joinLog.Role),
				string(joinLog.LogChannel),
			)
		}
	case serviceTerminationPosTemplateId:
		logPos := buffer.GetInt64(offset)
		adapter.agent.onServiceTerminationPosition(logPos)
	default:
		logger.Debugf("ServiceAdapter: unexpected templateId=%d at pos=%d", templateId, header.Position())
	}
}
