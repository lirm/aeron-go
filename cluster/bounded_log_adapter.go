package cluster

import (
	"bytes"
	"fmt"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type BoundedLogAdapter struct {
	marshaller     *codecs.SbeGoMarshaller
	options        *Options
	agent          *ClusteredServiceAgent
	image          *aeron.Image
	maxLogPosition int64
}

func (adapter *BoundedLogAdapter) IsDone() bool {
	return adapter.image.Position() >= adapter.maxLogPosition || adapter.image.IsEndOfStream() || adapter.image.IsClosed()
}

func (adapter *BoundedLogAdapter) Poll() int {
	return adapter.image.BoundedPoll(adapter.onFragment, adapter.maxLogPosition, adapter.options.LogFragmentLimit)
}

func (adapter *BoundedLogAdapter) onFragment(
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
			fmt.Println("join log decode error: ", err)
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
		fmt.Println("got termination position log")
	default:
		fmt.Println("unexpected template id: ", hdr.TemplateId)
	}
}

func (adapter *BoundedLogAdapter) Close() error {
	var err error
	if adapter.image != nil {
		err = adapter.image.Close()
		adapter.image = nil
	}
	return err
}
