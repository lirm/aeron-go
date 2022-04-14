package cluster

import (
	"bytes"
	"fmt"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ServiceAdapter struct {
	marshaller   *codecs.SbeGoMarshaller
	options      *Options
	agent        *ClusteredServiceAgent
	subscription *aeron.Subscription
}

// type FragmentHandler func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header)

func (adapter *ServiceAdapter) Poll() int {
	if adapter.subscription.IsClosed() {
		panic("subscription closed")
	}
	return adapter.subscription.Poll(adapter.fragmentAssembler, 10)
}

func (adapter *ServiceAdapter) fragmentAssembler(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	var hdr codecs.SbeGoMessageHeader
	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)

	if err := hdr.Decode(adapter.marshaller, buf); err != nil {
		fmt.Println("header decode error: ", err)
	}

	j := codecs.JoinLog{}
	t := codecs.ServiceTerminationPosition{}
	switch hdr.TemplateId {
	case j.SbeTemplateId():
		joinLog := &codecs.JoinLog{}
		if err := joinLog.Decode(
			adapter.marshaller,
			buf,
			hdr.Version,
			hdr.BlockLength,
			adapter.options.RangeChecking,
		); err != nil {
			fmt.Println("ServiceAdaptor: join log decode error: ", err)
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
	case t.SbeTemplateId():
		e := codecs.ServiceTerminationPosition{}
		if err := e.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("ServiceAdaptor: service termination pos decode error: ", err)
		} else {
			adapter.agent.onServiceTerminationPosition(e.LogPosition)
		}
	default:
		//fmt.Println("ServiceAdaptor: unexpected template id: ", hdr.TemplateId)
	}
}
