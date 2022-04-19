package cluster

import (
	"bytes"
	"fmt"

	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/cluster/codecs"
)

const (
	beginFrag    uint8 = 0x80
	endFrag      uint8 = 0x40
	unfragmented uint8 = 0x80 | 0x40
)

type BoundedLogAdapter struct {
	marshaller     *codecs.SbeGoMarshaller
	options        *Options
	agent          *ClusteredServiceAgent
	image          *aeron.Image
	builder        *bytes.Buffer
	maxLogPosition int64
}

func (adapter *BoundedLogAdapter) IsDone() bool {
	return adapter.image.Position() >= adapter.maxLogPosition ||
		adapter.image.IsEndOfStream() ||
		adapter.image.IsClosed()
}

func (adapter *BoundedLogAdapter) Poll(limitPos int64) int {
	return adapter.image.BoundedPoll(adapter.onFragment, limitPos, adapter.options.LogFragmentLimit)
}

func (adapter *BoundedLogAdapter) onFragment(
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

func (adapter *BoundedLogAdapter) onMessage(
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
		fmt.Printf("BoundedLogAdaptor - unexpected schemaId=%d templateId=%d blockLen=%d version=%d\n",
			schemaId, templateId, blockLength, version)
		return
	}
	offset += 8
	length -= 8

	switch templateId {
	case timerEventTemplateId:
		fmt.Println("BoundedLogAdaptor - got timer event")
	case sessionOpenTemplateId:
		event := &codecs.SessionOpenEvent{}
		if err := event.Decode(
			adapter.marshaller,
			toByteBuffer(buffer, offset, length),
			version,
			blockLength,
			adapter.options.RangeChecking,
		); err != nil {
			fmt.Println("session open decode error: ", err)
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
		event := &codecs.SessionCloseEvent{}
		if err := event.Decode(
			adapter.marshaller,
			toByteBuffer(buffer, offset, length),
			version,
			blockLength,
			adapter.options.RangeChecking,
		); err != nil {
			fmt.Println("session close decode error: ", err)
			return
		}

		adapter.agent.onSessionClose(
			event.LeadershipTermId,
			header.Position(),
			event.ClusterSessionId,
			event.Timestamp,
			event.CloseReason,
		)
	case clusterActionReqTemplateId:
		e := &codecs.ClusterActionRequest{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("cluster action request decode error: ", err)
		} else {
			adapter.agent.onServiceAction(e.LeadershipTermId, e.LogPosition, e.Timestamp, e.Action)
		}
	case newLeadershipTermTemplateId:
		e := &codecs.NewLeadershipTermEvent{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("new leadership term decode error: ", err)
		} else {
			//fmt.Println("BoundedLogAdaptor - got new leadership term: ", e)
			adapter.agent.onNewLeadershipTermEvent(e.LeadershipTermId, e.LogPosition, e.Timestamp, e.TermBaseLogPosition,
				e.LeaderMemberId, e.LogSessionId, e.TimeUnit, e.AppVersion)
		}
	case membershipChangeTemplateId:
		e := &codecs.MembershipChangeEvent{}
		buf := toByteBuffer(buffer, offset, length)
		if err := e.Decode(adapter.marshaller, buf, version, blockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("membership change event decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got membership change event: ", e)
		}
	case sessionMessageHeaderTemplateId:
		if length < SessionMessageHeaderLength {
			fmt.Println("received invalid session message - length: ", length)
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
		fmt.Println("BoundedLogAdaptor: unexpected template id: ", templateId)
	}
}

func toByteBuffer(buffer *atomic.Buffer, offset int32, length int32) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)
	return buf
}

func (adapter *BoundedLogAdapter) Close() error {
	var err error
	if adapter.image != nil {
		err = adapter.image.Close()
		adapter.image = nil
	}
	return err
}
