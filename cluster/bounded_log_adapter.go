package cluster

import (
	"bytes"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const (
	beginFrag    uint8 = 0x80
	endFrag      uint8 = 0x40
	unfragmented uint8 = 0x80 | 0x40
)

type boundedLogAdapter struct {
	marshaller     *codecs.SbeGoMarshaller
	options        *Options
	agent          *ClusteredServiceAgent
	image          *aeron.Image
	builder        *bytes.Buffer
	maxLogPosition int64
}

func (adapter *boundedLogAdapter) isDone() bool {
	return adapter.image.Position() >= adapter.maxLogPosition ||
		adapter.image.IsEndOfStream() ||
		adapter.image.IsClosed()
}

func (adapter *boundedLogAdapter) poll(limitPos int64) int {
	return adapter.image.BoundedPoll(adapter.onFragment, limitPos, adapter.options.LogFragmentLimit)
}

func (adapter *boundedLogAdapter) onFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	flags := header.Flags()
	if (flags & unfragmented) == unfragmented {
		adapter.onMessage(buffer, offset, length, header)
	} else if (flags & beginFrag) == beginFrag {
		if adapter.builder == nil {
			adapter.builder = &bytes.Buffer{}
		}
		adapter.builder.Reset()
		buffer.WriteBytes(adapter.builder, offset, length)
	} else if adapter.builder != nil && adapter.builder.Len() != 0 {
		buffer.WriteBytes(adapter.builder, offset, length)
		if (flags & endFrag) == endFrag {
			msgLength := adapter.builder.Len()
			adapter.onMessage(
				atomic.MakeBuffer(adapter.builder.Bytes(), msgLength),
				int32(0),
				int32(msgLength),
				header)
			adapter.builder.Reset()
		}
	}
}

func (adapter *boundedLogAdapter) onMessage(
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
	if schemaId != ClusterSchemaId {
		logger.Errorf("BoundedLogAdaptor - unexpected schemaId=%d templateId=%d blockLen=%d version=%d",
			schemaId, templateId, blockLength, version)
		return
	}
	offset += SBEHeaderLength
	length -= SBEHeaderLength

	switch templateId {
	case timerEventTemplateId:
		correlationId := buffer.GetInt64(offset + 8)
		timestamp := buffer.GetInt64(offset + 16)
		adapter.agent.onTimerEvent(header.Position(), correlationId, timestamp)
	case sessionOpenTemplateId:
		event := &codecs.SessionOpenEvent{}
		if err := event.Decode(
			adapter.marshaller,
			toByteBuffer(buffer, offset, length),
			version,
			blockLength,
			adapter.options.RangeChecking,
		); err != nil {
			logger.Errorf("boundedLogAdapter: session open decode error: %v", err)
			return
		}

		adapter.agent.onSessionOpen(
			event.LeadershipTermId,
			header.Position(),
			event.ClusterSessionId,
			event.Timestamp,
			event.ResponseStreamId,
			string(event.ResponseChannel),
			event.EncodedPrincipal,
		)
	case sessionCloseTemplateId:
		leadershipTermId := buffer.GetInt64(offset)
		clusterSessionId := buffer.GetInt64(offset + 8)
		timestamp := buffer.GetInt64(offset + 16)
		closeReason := codecs.CloseReasonEnum(buffer.GetInt32(offset + 24))
		adapter.agent.onSessionClose(leadershipTermId, header.Position(), clusterSessionId, timestamp, closeReason)
	case clusterActionReqTemplateId:
		e := codecs.ClusterActionRequest{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			logger.Errorf("boundedLogAdapter: cluster action request decode error: %v", err)
		} else {
			adapter.agent.onServiceAction(e.LeadershipTermId, e.LogPosition, e.Timestamp, e.Action)
		}
	case newLeadershipTermTemplateId:
		e := codecs.NewLeadershipTermEvent{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			logger.Errorf("boundedLogAdapter: new leadership term decode error: %v", err)
		} else {
			adapter.agent.onNewLeadershipTermEvent(e.LeadershipTermId, e.LogPosition, e.Timestamp,
				e.TermBaseLogPosition, e.LeaderMemberId, e.LogSessionId, e.TimeUnit, e.AppVersion)
		}
	case membershipChangeTemplateId:
		e := codecs.MembershipChangeEvent{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			logger.Errorf("boundedLogAdapter: membership change event decode error: %v", err)
		} else {
			adapter.agent.onMembershipChange(e.LogPosition, e.Timestamp, e.ChangeType, e.MemberId)
		}
	case SessionMessageHeaderTemplateId:
		if length < SessionMessageHeaderLength {
			logger.Errorf("received invalid session message - length: %d", length)
			return
		}
		clusterSessionId := buffer.GetInt64(offset + 8)
		timestamp := buffer.GetInt64(offset + 16)
		adapter.agent.onSessionMessage(
			header.Position(),
			clusterSessionId,
			timestamp,
			buffer,
			offset+SessionMessageHeaderLength,
			length-SessionMessageHeaderLength,
			header,
		)
	default:
		logger.Debugf("BoundedLogAdaptor: unexpected templateId=%d at pos=%d", templateId, header.Position())
	}
}

func toByteBuffer(buffer *atomic.Buffer, offset int32, length int32) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)
	return buf
}

func (adapter *boundedLogAdapter) Close() error {
	var err error
	if adapter.image != nil {
		err = adapter.image.Close()
		adapter.image = nil
	}
	return err
}
